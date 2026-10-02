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

type UploaderEncryptApplication interface {
	UploadEncrypt(IncomeData) (string, error)
}

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
		return "", ErrorFileSizeBig
	}
	g, ctx := errgroup.WithContext(data.Ctx)
	x, err := sa.Uploader.SetName(data.Name).SetSize(data.Size).Make(ctx)
	if err != nil {
		return "", err
	}
	shortNameFile := sa.Generate.GenerateText(4)
	allData := memguard.NewBuffer(sa.Encrypt.GetRequiredOverheadSize() + 32)
	defer allData.Destroy()
	err = fillOut(allData.Data(), sa.Encrypt.GetRequiredOverheadSize())
	if err != nil {
		return "", err
	}
	g.Go(func() error {
		err := sa.EncryptFile(data.File, x, allData.Data())
		if err != nil {
			return err
		}
		return nil
	})

	FileInfoInBytes, err := sa.Encode.Encode(DomainLevel.FileLabelsBytes{
		FileName: data.Name,
		AesKey:   hex.EncodeToString(allData.Data()[aes.BlockSize:]),
	})
	if err != nil {
		slog.Error("UploadEncrypt; error to encode data into Json", "ERROR", err)
		return "", errors.New(ErrorUploadingEncrypt)
	}
	encr, err := sa.Encrypt.MakeCrypto(sa.ServerKeys.GerOurPrivateKey(), nil)
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

func (sa *NewUploadEncrypt) EncryptFile(src io.Reader, dst DomainLevel.Upload, data []byte) error {
	cryptoData, err := sa.Encrypt.MakeCrypto(data[sa.Encrypt.GetRequiredOverheadSize():], data[:sa.Encrypt.GetRequiredOverheadSize()])
	if err != nil {
		return err
	}
	buf := memguard.NewBuffer(32 * 1024)
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
