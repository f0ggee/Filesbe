package Application

import (
	"Kaban/internal/DomainLevel"

	"context"
	"errors"
	"io"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

const (
	FileMaxSize = 500000000
)

var (
	ErrorFileSizeBig = errors.New("the File's size is bigger than the default size")
)

type UploadApplication interface {
	Upload(FileUploaderIncomeData) (string, error)
}

type NewFileUploaderDataMange struct {
	Encode DomainLevel.Encoder
}
type NewFileUploaderCrypto struct {
	Generator DomainLevel.CryptoGenerating
}

type NewFileUploaderDelivery struct {
	Uploader   DomainLevel.MakerUploader
	WriteRedis DomainLevel.WritingRedis
}
type NewUpload struct {
	NewFileUploaderCrypto
	NewFileUploaderDataMange
	NewFileUploaderDelivery
}

func GetNewFileUploader(newFileUploaderCrypto NewFileUploaderCrypto, newFileUploaderDataMange NewFileUploaderDataMange, newFileUploaderDelivery NewFileUploaderDelivery) *NewUpload {
	return &NewUpload{NewFileUploaderCrypto: newFileUploaderCrypto, NewFileUploaderDataMange: newFileUploaderDataMange, NewFileUploaderDelivery: newFileUploaderDelivery}
}

type FileUploaderIncomeData struct {
	File io.ReadCloser
	Name string
	Size int64
	Ctx  context.Context
}

func (sa *NewUpload) Upload(r FileUploaderIncomeData) (string, error) {
	g, ctx := errgroup.WithContext(r.Ctx)
	if r.Size >= FileMaxSize {
		return "", ErrorFileSizeBig
	}
	defer func() {
		err := r.File.Close()
		if err != nil {
			slog.Error("Upload; error to close a body", "ERROR", err)
			return
		}
	}()
	shortNameFile := sa.Generator.GenerateText(4)

	g.Go(func() error {
		uploader, err := sa.Uploader.SetName(r.Name).SetSize(r.Size).Make(ctx)
		if err != nil {
			return err
		}
		defer func(uploader DomainLevel.Upload) {
			err := uploader.CloseSource()
			if err != nil {
				return
			}
		}(uploader)
		err = uploader.Uploader(r.File)
		if err != nil {
			return err
		}
		return nil
	})

	fileIntoBytes, err := sa.Encode.Encode(r.Name)
	if err != nil {
		return "", err
	}

	if err = g.Wait(); err != nil {
		return "", err
	}
	err = sa.WriteRedis.WriteData(DomainLevel.WriteDataIncomeData{
		FileName: shortNameFile,
		Info:     fileIntoBytes,
		Ctx:      r.Ctx,
	})
	if err != nil {
		return "", err
	}
	return shortNameFile, nil

}
