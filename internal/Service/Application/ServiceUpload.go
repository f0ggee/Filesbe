package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"Kaban/internal/InfrastructureLayer/s3Repo"
	"context"
	"errors"
	"io"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

const (
	FileMaxSize         = 500000000
	ErrorFileSizeBig    = "the File's size is bigger than the default size"
	ErrorStartUploading = "an unexpected error happened"
)

type NewFileUploaderDataMange struct {
	Encode RepoParsers.Encode
}
type NewFileUploaderCrypto struct {
	Generator DomainLevel.CryptoGenerating
}

type NewFileUploaderDelivery struct {
	UploadS3   s3Repo.S3Uploader
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

func (sa *NewUpload) FileUploader(r FileUploaderIncomeData) (string, error) {
	g, ctx := errgroup.WithContext(r.Ctx)
	if r.Size >= FileMaxSize {
		return "", errors.New(ErrorFileSizeBig)
	}
	defer func() {
		err := r.File.Close()
		if err != nil {
			slog.Error("FileUploader; error to close a body", "ERROR", err)
			return
		}
	}()
	shortNameFile := sa.Generator.GenerateText(4)
	settings := DomainLevel.GetNewFileSettings(r.Size, r.Name)
	Parts, goroutines := settings.FindBestOptions()
	fileFormat := settings.FindFormatOfFile()

	g.Go(func() error {
		err2 := sa.UploadS3.UploadFile(s3Repo.UploadFileIncomingData{
			Parts:      Parts,
			Goroutines: goroutines,
			Ctx:        ctx,
			FileDetails: s3Repo.FileDetails{
				FileFormat: fileFormat,
				FileName:   r.Name,
				FileBody:   s3Repo.TypeUploading{Normal: r.File},
			},
		})
		if err2 != nil {
			return err2
		}
		return nil
	})

	fileIntoBytes, err := sa.Encode.JsonEncodeMarshall(r.Name)
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
