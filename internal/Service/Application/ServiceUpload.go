package Application

import (
	"Kaban/internal/DomainLevel"
	"errors"
	"log/slog"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type NewFileUploader struct {
	GetCrypto
	GetFileManager
	S3Controlling
	Parser
	RedisControlling
}

func GetNewNewFileUploader(getCrypto GetCrypto, getFileManager GetFileManager, s3Controlling S3Controlling, parser Parser, redisControlling RedisControlling) *NewFileUploader {
	return &NewFileUploader{GetCrypto: getCrypto, GetFileManager: getFileManager, S3Controlling: s3Controlling, Parser: parser, RedisControlling: redisControlling}
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
		err2 := sa.Uploader.UploadFile(DomainLevel.UploadFileIncomingData{
			Parts:      Parts,
			Goroutines: goroutines,
			Ctx:        ctx,
			FileDetails: DomainLevel.FileDetails{
				FileFormat: fileFormat,
				FileName:   fileDetails.Filename,
				FileBody:   DomainLevel.TypeUploading{Normal: file},
			},
		})
		if err2 != nil {
			return err2
		}
		return nil
	})

	fileIntoBytes, err := sa.Parser.Encode.JsonEncodeMarshall(fileDetails.Filename)
	if err != nil {
		return "", err
	}
	if err := g.Wait(); err != nil {
		return "", err
	}
	err = sa.RedisControlling.Writer.WriteData(shortNameFile, fileIntoBytes, r.Context())
	if err != nil {
		return "", err
	}
	return shortNameFile, nil

}
