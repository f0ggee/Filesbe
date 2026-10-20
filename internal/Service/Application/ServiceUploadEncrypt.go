package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"bytes"

	"context"
	"crypto/rand"
	"errors"
	"io"
	"log/slog"
	"time"

	"github.com/awnumar/memguard"
	"golang.org/x/sync/errgroup"
)

const ErrorUploadingEncrypt = "an unexpected error occurred while loading encryption"

type NewUploadEncryptCrypto struct {
	Generate      DomainLevel.CryptoGenerating
	Encrypt       DomainLevel.CryptoMaker
	Rsa           DomainLevel.CryptoMaker
	EncrypterKeys RepoEncrypterKeys.Keys
}
type NewUploadEncryptDataManage struct {
	Encode DomainLevel.Encoder
}
type NewUploadEncryptDelivery struct {
	Uploader     DomainLevel.MakerUploader
	RedisWriter  DomainLevel.WritingRedis
	RedisChecker DomainLevel.RedisChecker
	RedisDeleter DomainLevel.DeleterRedis
	Deleter      DomainLevel.MakerDeleter
}
type UploadEncrypt struct {
	NewUploadEncryptDataManage
	NewUploadEncryptCrypto
	NewUploadEncryptDelivery
}

func GetNewUploadEncrypt(newUploadEncryptDataManage NewUploadEncryptDataManage, newUploadEncryptCrypto NewUploadEncryptCrypto, newUploadEncryptDelivery NewUploadEncryptDelivery) *UploadEncrypt {
	return &UploadEncrypt{NewUploadEncryptDataManage: newUploadEncryptDataManage, NewUploadEncryptCrypto: newUploadEncryptCrypto, NewUploadEncryptDelivery: newUploadEncryptDelivery}
}

type UploadEncryptIncomeData struct {
	File io.Reader
	Name string
	Ctx  context.Context
	Size int64
}

func (sa *UploadEncrypt) UploadEncryptApplication(data UploadEncryptIncomeData) (string, error) {
	if data.Size >= FileMaxSize {
		return "", ErrorFileSizeBig
	}

	g, ctx := errgroup.WithContext(data.Ctx)

	x, err := sa.Uploader.SetName(data.Name).SetSize(data.Size).Make(ctx)
	if err != nil {
		return "", err
	}

	defer x.CloseSource()

	shortNameFile := sa.Generate.GenerateText(4)
	allData := memguard.NewBuffer(sa.Encrypt.GetRequiredOverheadSize() + 32)
	defer allData.Destroy()
	fillOut(allData.Data())

	g.Go(func() error {
		err := sa.encryptFile(data.File, x, allData.Data())
		if err != nil {
			return err
		}
		return nil
	})

	FileInfoInBytes, err := sa.Encode.Encode(DomainLevel.FileLabelsBytes{
		FileName: data.Name,
		AesKey:   allData.Data()[sa.Encrypt.GetRequiredOverheadSize():],
	})
	if err != nil {
		slog.Error("UploadEncrypt; error to encode data into Json", "ERROR", err)
		return "", errors.New(ErrorUploadingEncrypt)
	}
	encr, err := sa.Rsa.MakeCrypto(sa.EncrypterKeys.GetKey(), nil)
	if err != nil {
		return "", err
	}
	encryptedFileInfo, err := encr.Encrypt(FileInfoInBytes)
	if err != nil {
		return "", err
	}
	if err = g.Wait(); err != nil {
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

func fillOut(allData []byte) {
	if _, err := io.ReadFull(rand.Reader, allData); err != nil {
		slog.Error("EncryptFile: error to generate a aesKey", "ERROR", err)
		return
	}
	return
}

func (sa *UploadEncrypt) encryptFile(src io.Reader, dst DomainLevel.Upload, data []byte) error {
	cryptoData, err := sa.Encrypt.MakeCrypto(data[sa.Encrypt.GetRequiredOverheadSize():], data[:sa.Encrypt.GetRequiredOverheadSize()])
	if err != nil {
		return err
	}
	buf := memguard.NewBuffer(DomainLevel.ChunkSize)
	defer buf.Destroy()
	for {
		n, err := src.Read(buf.Bytes())
		if n > 0 {
			outData, err := cryptoData.Encrypt(buf.Bytes()[:n])
			if err != nil {
				return err
			}
			err = dst.Uploader(bytes.NewReader(outData))
			if err != nil {
				slog.Error("EncryptFile; error to write in the stream", "ERROR", err.Error())
				return err
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

func (sa *UploadEncrypt) setDeleteFile(shortNameFile string) *time.Timer {
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
