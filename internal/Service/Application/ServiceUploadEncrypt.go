package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/Crypto"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"Kaban/internal/InfrastructureLayer/s3Repo"
	"context"
	"crypto/aes"
	"crypto/cipher"
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
	Generate   DomainLevel.CryptoGenerating
	ServerKeys DomainLevel.NewServerKeys
	Encrypt    DomainLevel.Encrypt
}
type NewUploadEncryptDataManage struct {
	Encode RepoParsers.Encode
}
type NewUploadEncryptDelivery struct {
	UploaderS3   s3Repo.S3Uploader
	RedisWriter  DomainLevel.WritingRedis
	RedisChecker DomainLevel.RedisChecker
	RedisDeleter DomainLevel.DeleterRedis
	DeleterS3    s3Repo.DeleterS3
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
	settings := DomainLevel.GetNewFileSettings(data.Size, data.Name)
	BesParts, goroutine, FileExtension := sa.getFileSettings(settings)
	shortNameFile := sa.Generate.GenerateText(4)

	EncryptAesKey, err := memguard.NewBufferFromReader(rand.Reader, 32)
	if err != nil {
		slog.Error("UploadEncrypt: error to generate random bytes", "ERROR", err)
		return "", errors.New(ErrorUploadingEncrypt)
	}
	g.Go(func() error {
		err = sa.EncryptFile(data.File, writer, EncryptAesKey.Data())
		if err != nil {
			return err
		}
		return nil
	})
	defer EncryptAesKey.Destroy()

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
		AesKey:   hex.EncodeToString(EncryptAesKey.Bytes()),
	})
	if err != nil {
		slog.Error("UploadEncrypt; error to encode data into Json", "ERROR", err)
		return "", errors.New(ErrorUploadingEncrypt)
	}

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		errS3 := sa.UploaderS3.UploadFileEncrypt(s3Repo.UploadFileIncomingData{
			Parts:      BesParts,
			Goroutines: goroutine,
			Ctx:        ctx,
			FileDetails: s3Repo.FileDetails{
				FileFormat: FileExtension,
				FileName:   shortNameFile,
				FileBody: s3Repo.TypeUploading{
					Pipe: reader,
				},
			},
		})
		if errS3 != nil {
			return errS3
		}
		return nil
	})
	EncryptFileInfo, err := sa.Encrypt.Encrypter(DomainLevel.IncomeEncryptData{
		Key:  x509.MarshalPKCS1PublicKey(Public.Public().(*rsa.PublicKey)),
		Data: FileInfoInBytes,
	})
	if err != nil {
		return "", err
	}
	err = sa.RedisWriter.WriteData(DomainLevel.WriteDataIncomeData{
		FileName: shortNameFile,
		Info:     EncryptFileInfo,
		Ctx:      data.Ctx,
	})
	if err != nil {
		err := writer.CloseWithError(err)
		return "", err
	}

	sa.setDeleteFile(shortNameFile)

	return shortNameFile, nil

}

func (sa *NewUploadEncrypt) getFileSettings(settings *DomainLevel.FileSettings) (int, int, string) {
	BesParts, goroutine := settings.FindBestOptions()
	FileExtension := settings.FindFormatOfFile()
	return BesParts, goroutine, FileExtension
}

func (sa *NewUploadEncrypt) EncryptFile(file io.ReadCloser, writer io.Writer, aesKey []byte) error {

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		slog.Error("EncryptFile; error to make a new cipher block", "ERROR", err)
		return err
	}

	nonce := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	stream := cipher.NewCTR(block, nonce)
	buf := memguard.NewBuffer(32 * 1024)
	defer buf.Destroy()
	_, err = writer.Write(nonce)
	if err != nil {
		slog.Error("EncryptFile; error happened during writing", "ERROR", err)
		return errors.New(Crypto.ErrorEncryptFile)
	}
	for {
		n, err := file.Read(buf.Bytes())
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("EncryptFile; error to download a File", "ERROR", err.Error())
			return errors.New(ErrorUploadingEncrypt)
		}
		stream.XORKeyStream(buf.Bytes()[:n], buf.Bytes()[:n])
		_, err = writer.Write(buf.Bytes()[:n])
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
			err := sa.DeleterS3.DeleteFileFromS3(shortNameFile, Ctx)
			if err != nil {
				return err
			}
			return nil
		})
		if err := g2.Wait(); err != nil {
			slog.Error("DeleterS3; error to delete a File", "ERROR", err)
			return
		}
		return
	})
}
