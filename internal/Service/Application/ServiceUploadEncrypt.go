package Application

import (
	"Kaban/internal/DomainLevel"
	"bytes"

	"context"
	"crypto/aes"
	"crypto/rand"
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
	Encrypt    DomainLevel.CryptoMaker
}
type NewUploadEncryptDataManage struct {
	Encode DomainLevel.Encoder
}
type NewUploadEncryptDelivery struct {
	Uploader     DomainLevel.MakerUploader
	RedisWriter  DomainLevel.WritingRedis
	RedisChecker DomainLevel.RedisChecker
	RedisDeleter DomainLevel.DeleterRedis

	Deleter DomainLevel.MakerDeleter
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
	Ctx  context.Context
	Size int64
}

func (sa *NewUploadEncrypt) UploadEncrypt(data IncomeData) (string, error) {
	if data.Size >= FileMaxSize {
		return "", errors.New(ErrorFileSizeBig)
	}
	reader, writer := io.Pipe()
	defer func() {
		closeSources(data.File, reader, writer)
	}()

	g, ctx := errgroup.WithContext(data.Ctx)
	x, err := sa.Uploader.SetName(data.Name).SetSize(data.Size).Make(ctx)
	if err != nil {
		return "", err
	}
	shortNameFile := sa.Generate.GenerateText(4)
	allData := memguard.NewBuffer(sa.Encrypt.GetRequiredRandomSize() + 32)
	err = fillOut(allData.Data(), sa.Encrypt.GetRequiredRandomSize())
	if err != nil {
		return "", err
	}
	sa.setEncrypter(data, g, x, allData)
	defer allData.Destroy()

	FileInfoInBytes, err := sa.Encode.Encode(DomainLevel.FileLabelsBytes{
		FileName: data.Name,
		AesKey:   hex.EncodeToString(allData.Data()[aes.BlockSize:]),
	})
	if err != nil {
		slog.Error("UploadEncrypt; error to encode data into Json", "ERROR", err)
		return "", errors.New(ErrorUploadingEncrypt)
	}

	setUploader(g, x, reader)
	encr, err := sa.Encrypt.MakeCrypto(sa.ServerKeys.GerOurPrivateKey(), []byte("1"))
	if err != nil {
		return "", err
	}
	encryptedFileInfo, err := encr.Encrypt(FileInfoInBytes)
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

func setUploader(g *errgroup.Group, readyUploader DomainLevel.Upload, src io.Reader) {
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

func (sa *NewUploadEncrypt) setEncrypter(data IncomeData, g *errgroup.Group, d DomainLevel.Upload, allData *memguard.LockedBuffer) {
	g.Go(func() error {
		err := sa.EncryptFile(data.File, d, allData.Data())
		if err != nil {
			return err
		}
		return nil
	})
}

func fillOut(allData []byte, overHeadSize int) error {
	if _, err := io.ReadFull(rand.Reader, allData[:overHeadSize]); err != nil {
		slog.Error("EncryptFile: error to generate a nonce", "ERROR", err)
		return errors.New(ErrorUploadingEncrypt)
	}
	if _, err := io.ReadFull(rand.Reader, allData[overHeadSize:]); err != nil {
		slog.Error("EncryptFile: error to generate a aesKey", "ERROR", err)
		return errors.New(ErrorUploadingEncrypt)
	}
	return nil
}

func (sa *NewUploadEncrypt) EncryptFile(file io.Reader, d DomainLevel.Upload, data []byte) error {
	cryptoData, err := sa.Encrypt.MakeCrypto(data[sa.Encrypt.GetRequiredRandomSize():], data[:sa.Encrypt.GetRequiredRandomSize()])
	if err != nil {
		return err
	}
	buf := memguard.NewBuffer(32 * 1024)
	defer buf.Destroy()
	for {
		n, err := file.Read(buf.Bytes())
		if n > 0 {
			outData, err := cryptoData.Encrypt(buf.Bytes()[:n])
			if err != nil {
				return err
			}
			err = d.Uploader(bytes.NewReader(outData))
			if err != nil {
				slog.Error("EncryptFile; error to write in the stream", "ERROR", err.Error())
				return errors.New(DomainLevel.ErrorStrangeCrypto)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("EncryptFile; error to download a File", "ERROR", err.Error())
			return errors.New(ErrorUploadingEncrypt)
		}

	}

	return nil
}

func closeSources(data io.ReadCloser, reader *io.PipeReader, writer *io.PipeWriter) {
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
