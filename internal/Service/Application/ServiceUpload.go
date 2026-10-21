package Application

import (
	"Kaban/internal/DomainLevel"

	"context"
	"io"

	"golang.org/x/sync/errgroup"
)

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

type UploadIncomeData struct {
	File io.Reader
	Name string
	Size int64
	Ctx  context.Context
}

func (sa *NewUpload) UploadApplication(r UploadIncomeData) (string, error) {
	g, ctx := errgroup.WithContext(r.Ctx)
	if r.Size >= FileMaxSize {
		return "", ErrorFileSizeBig
	}

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
		err = uploader.Upload(r.File)
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
