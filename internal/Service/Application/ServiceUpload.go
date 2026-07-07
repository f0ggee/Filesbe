package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/s3Repo"
	"errors"
	"log/slog"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type NewFileUploader struct {
	getCrypto
	getFileManager
	s3Controlling
	parser
	redisControlling
}

func GetNewNewFileUploader(getCrypto getCrypto, getFileManager getFileManager, s3Controlling s3Controlling, parser parser, redisControlling redisControlling) *NewFileUploader {
	return &NewFileUploader{getCrypto: getCrypto, getFileManager: getFileManager, s3Controlling: s3Controlling, parser: parser, redisControlling: redisControlling}
}

func (sa *NewFileUploader) FileUploader(r *http.Request) (string, error) {
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

	shortNameFile := sa.Generate.GenerateShortName()

	Parts, goroutines := sa.FileManaging.FindBestOptions(fileDetails.Size)

	fileFormat := sa.FileManaging.FindFormatOfFile(fileDetails.Filename)
	g.Go(func() error {
		err2 := sa.Uploader.UploadFile(s3Repo.UploadFileIncomingData{
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

	fileIntoBytes, err := sa.parser.Encode.JsonEncodeMarshall(fileDetails.Filename)
	if err != nil {
		return "", err
	}
	if err := g.Wait(); err != nil {
		return "", err
	}
	err = sa.redisControlling.Writer.WriteData(shortNameFile, fileIntoBytes, r.Context())
	if err != nil {
		return "", err
	}
	return shortNameFile, nil

}
