package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/FileControls"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"Kaban/internal/InfrastructureLayer/s3Repo"
	"errors"
	"log/slog"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type NewFileUploaderDataMange struct {
	FileSettings FileControls.FileSettings
	Encode       RepoParsers.Encode
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

func (sa *NewUpload) FileUploader(r *http.Request) (string, error) {
	g, ctx := errgroup.WithContext(r.Context())
	file, fileDetails, err := r.FormFile("file")
	if err != nil {
		slog.Error("FileUploader; error to get a file", "ERROR", err)
		return "", errors.New(DomainLevel.ErrorStrangeUploadFile)
	}
	if fileDetails.Size >= DomainLevel.FileMaxSize {
		return "", errors.New(DomainLevel.ErrorFileSizeBig)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			slog.Error("FileUploader; error to close a body", "ERROR", err)
			return
		}
	}()

	shortNameFile := sa.Generator.GenerateShortName()

	Parts, goroutines := sa.FileSettings.FindBestOptions(fileDetails.Size)

	fileFormat := sa.FileSettings.FindFormatOfFile(fileDetails.Filename)
	g.Go(func() error {
		err2 := sa.UploadS3.UploadFile(s3Repo.UploadFileIncomingData{
			Parts:      Parts,
			Goroutines: goroutines,
			Ctx:        ctx,
			FileDetails: s3Repo.FileDetails{
				FileFormat: fileFormat,
				FileName:   fileDetails.Filename,
				FileBody:   s3Repo.TypeUploading{Normal: file},
			},
		})
		if err2 != nil {
			return err2
		}
		return nil
	})

	fileIntoBytes, err := sa.Encode.JsonEncodeMarshall(fileDetails.Filename)
	if err != nil {
		return "", err
	}
	if err := g.Wait(); err != nil {
		return "", err
	}
	err = sa.WriteRedis.WriteData(DomainLevel.WriteDataIncomeData{
		FileName: shortNameFile,
		Info:     fileIntoBytes,
		Ctx:      r.Context(),
	})
	if err != nil {
		return "", err
	}
	return shortNameFile, nil

}
