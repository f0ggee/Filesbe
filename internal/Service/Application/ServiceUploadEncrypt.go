package Application

import (
	"Kaban/internal/DomainLevel"
	s3Repo2 "Kaban/internal/InfrastructureLayer/FileTransferring/s3Repo"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"context"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"time"

	"github.com/awnumar/memguard"
	"golang.org/x/sync/errgroup"
)

const ErrorUploadingEncrypt = "an unexpected error occurred while loading encryption"

type NewUploadEncryptCrypto struct {
	Generate          DomainLevel.CryptoGenerating
	ServerKeys        DomainLevel.NewServerKeys
	Encrypt           DomainLevel.CryptoMaker
	EncryptSecondMode DomainLevel.CryptoMakerSpec
}
type NewUploadEncryptDataManage struct {
	Encode RepoParsers.Encode
}
type NewUploadEncryptDelivery struct {
	UploaderS3   s3Repo2.S3Uploader
	Uploader     DomainLevel.MakerUploader
	RedisWriter  DomainLevel.WritingRedis
	RedisChecker DomainLevel.RedisChecker
	RedisDeleter DomainLevel.DeleterRedis
	DeleterS3    s3Repo2.DeleterS3
	Deleter      DomainLevel.MakerDeleter
}
type NewUploadEncrypt struct {
	NewUploadEncryptDataManage
	NewUploadEncryptCrypto
	NewUploadEncryptDelivery
}

func GetNewUploadEncrypt(newUploadEncryptDataManage NewUploadEncryptDataManage, newUploadEncryptCrypto NewUploadEncryptCrypto, newUploadEncryptDelivery NewUploadEncryptDelivery) *NewUploadEncrypt {
	return &NewUploadEncrypt{NewUploadEncryptDataManage: newUploadEncryptDataManage, NewUploadEncryptCrypto: newUploadEncryptCrypto, NewUploadEncryptDelivery: newUploadEncryptDelivery}
}

type IncomeData struct {
	File io.ReadCloser
	Name string
	Size int64
	Ctx  context.Context
}

func (sa *NewUploadEncrypt) UploadEncrypt(data IncomeData) (string, error) {
	if data.Size >= FileMaxSize {
		return "", errors.New(ErrorFileSizeBig)
	}
	reader, writer := io.Pipe()
	defer func() {
		sa.closeSources(data.File, reader, writer)
	}()

	g, ctx := errgroup.WithContext(data.Ctx)

	shortNameFile := sa.Generate.GenerateText(4)
	allData := memguard.NewBuffer(aes.BlockSize + 32)
	err := sa.fillOut(allData.Data())
	if err != nil {
		return "", err
	}
	sa.setEncrypter(data, g, writer, allData)
	defer allData.Destroy()

	Public, err := x509.ParsePKCS1PrivateKey(sa.ServerKeys.GerOurPrivateKey())
	if err != nil {
		slog.Error("UploadEncrypt; error to decode a key", "ERROR", err)
		if err != nil {
			slog.Error("Error in File writing wile a key parsing ", "Error", err)
			return "", errors.New(ErrorUploadingEncrypt)
		}
		return "", errors.New(ErrorUploadingEncrypt)
	}

	FileInfoInBytes, err := sa.Encode.JsonEncodeMarshall(DomainLevel.FileLabelsBytes{
		FileName: data.Name,
		AesKey:   hex.EncodeToString(allData.Data()[aes.BlockSize:]),
	})
	if err != nil {
		slog.Error("UploadEncrypt; error to encode data into Json", "ERROR", err)
		return "", errors.New(ErrorUploadingEncrypt)
	}

	x, err := sa.Uploader.SetName(data.Name).SetSize(data.Size).Make(ctx)
	if err != nil {
		return "", err
	}
	sa.setUploader(g, x, reader)
	key := x509.MarshalPKCS1PublicKey(Public.Public().(*rsa.PublicKey))
	encr, err := sa.EncryptSecondMode.MakeCrypto(key, 1)
	if err != nil {
		return "", err
	}
	encryptedFileInfo, err := encr.Encrypter(FileInfoInBytes)
	if err != nil {
		return "", err
	}
	err = sa.RedisWriter.WriteData(DomainLevel.WriteDataIncomeData{
		FileName: shortNameFile,
		Info:     encryptedFileInfo,
		Ctx:      data.Ctx,
	})
	if err != nil {
		return "", err
	}
	sa.setDeleteFile(shortNameFile)

	return shortNameFile, nil

}

func (sa *NewUploadEncrypt) setUploader(g *errgroup.Group, readyUploader DomainLevel.Upload, src io.Reader) {
	g.Go(func() error {
		defer func(readyUploader DomainLevel.Upload) {
			err := readyUploader.CloseSource()
			if err != nil {
				return
			}
		}(readyUploader)
		if err := readyUploader.Uploader(src); err != nil {
			return err
		}
		return nil
	})
}

func (sa *NewUploadEncrypt) setEncrypter(data IncomeData, g *errgroup.Group, writer *io.PipeWriter, allData *memguard.LockedBuffer) {
	g.Go(func() error {
		err := sa.EncryptFile(data.File, writer, allData.Data())
		if err != nil {
			return err
		}
		return nil
	})
}

func (sa *NewUploadEncrypt) fillOut(allData []byte) error {
	if _, err := io.ReadFull(rand.Reader, allData[:aes.BlockSize]); err != nil {
		slog.Error("EncryptFile: error to generate a nonce", "ERROR", err)
		return errors.New(ErrorUploadingEncrypt)
	}
	if _, err := io.ReadFull(rand.Reader, allData[aes.BlockSize:]); err != nil {
		slog.Error("EncryptFile: error to generate a aesKey", "ERROR", err)
		return errors.New(ErrorUploadingEncrypt)
	}
	return nil
}

func (sa *NewUploadEncrypt) getFileSettings(settings *DomainLevel.FileSettings) (int, int, string) {
	BesParts, goroutine := settings.FindBestOptions()
	FileExtension := settings.FindFormatOfFile()
	return BesParts, goroutine, FileExtension
}

func (sa *NewUploadEncrypt) EncryptFile(file io.Reader, writer io.Writer, data []byte) error {
	cryptoData, err := sa.Encrypt.MakeCrypto(data)
	if err != nil {
		return err
	}
	buf := memguard.NewBuffer(32 * 1024)
	defer buf.Destroy()
	for {
		n, err := file.Read(buf.Bytes())
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("EncryptFile; error to download a File", "ERROR", err.Error())
			return errors.New(ErrorUploadingEncrypt)
		}
		outData, err := cryptoData.Encrypter(buf.Bytes()[:n])
		if err != nil {
			return err
		}
		_, err = writer.Write(outData)
		if err != nil {
			slog.Error("EncryptFile; error to write in the stream ", "ERROR", err.Error())
			return errors.New(DomainLevel.ErrorStrangeCrypto)
		}
	}

	return nil
}

func (sa *NewUploadEncrypt) closeSources(data io.ReadCloser, reader *io.PipeReader, writer *io.PipeWriter) {
	err := data.Close()
	if err != nil {
		slog.Error("UploadEncrypt; error to close a File flow", "ERROR", err)
		return
	}
	err = reader.Close()
	if err != nil {
		slog.Error("UploadEncrypt: error close a reader", "ERROR", err)
		return
	}
	err = writer.Close()
	if err != nil {
		slog.Error("UploadEncrypt: error to close a writer", "ERROR", err)
	}
}

func (sa *NewUploadEncrypt) setDeleteFile(shortNameFile string) *time.Timer {
	return time.AfterFunc(5*time.Minute, func() {
		g2, Ctx := errgroup.WithContext(context.Background())
		Ctx, cancel := context.WithTimeout(Ctx, 25*time.Second)
		defer cancel()
		DownloadingHaveStarted := sa.RedisChecker.ChekIsStartDownload(shortNameFile, Ctx)
		if DownloadingHaveStarted {
			return
		}
		g2.Go(func() error {

			err := sa.RedisDeleter.DeleteFileInfo(shortNameFile, Ctx)
			if err != nil {
				return err
			}
			return nil
		})
		g2.Go(func() error {
			s, err := sa.Deleter.SetName(shortNameFile).Make(Ctx)
			if err != nil {
				return err
			}
			return s.Deleter()
		})
		if err := g2.Wait(); err != nil {
			slog.Error("DeleterS3; error to delete a File", "ERROR", err)
			return
		}
		return
	})
}
