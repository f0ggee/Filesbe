package s3Repo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
)

type NewUploading struct {
	S3Info Variables
}

func GetNewUploading(s3Info Variables) *NewUploading {
	return &NewUploading{S3Info: s3Info}
}

type TypeUploading struct {
	Pipe   *io.PipeReader
	Normal io.ReadCloser
}

type FileDetails struct {
	FileFormat string
	FileName   string
	FileBody   TypeUploading
}
type UploadFileIncomingData struct {
	Parts      int
	Goroutines int
	Ctx        context.Context
	FileDetails
}

type S3Uploader interface {
	UploadFile(UploadFileIncomingData) error

	UploadFileEncrypt(UploadFileIncomingData) error
}

func (sa *NewUploading) UploadFileEncrypt(data UploadFileIncomingData) error {
	slog.Group("File uploading details",
		slog.String("FileExtension", data.FileFormat),
		slog.String("Parts", fmt.Sprint(data.Parts)),
		slog.String("Goroutines", fmt.Sprint(data.Goroutines)),
	)

	uploader := manager.NewUploader(sa.S3Info.S3Connect, func(uploader *manager.Uploader) {

		uploader.MaxUploadParts = 200
		uploader.PartSize = int64(data.Parts) * 1024 * 1024
		uploader.Concurrency = data.Goroutines
		uploader.BufferProvider = manager.NewBufferedReadSeekerWriteToPool(data.Parts)
	})

	_, err := uploader.Upload(data.Ctx, &s3.PutObjectInput{
		Bucket:      aws.String(sa.S3Info.Bucket),
		Key:         aws.String(data.FileName),
		ContentType: aws.String(data.FileFormat),
		Body:        data.FileBody.Pipe,
	})
	if err == nil {
		return nil
	}

	var ns *types.NoSuchKey

	switch {

	case errors.As(err, &ns):

		slog.Error("file was used", "Error", err.Error())
		return errors.New(ErrorCantFindFile)

	case errors.Is(err, context.Canceled):
		slog.Error("file downloading was cancelled")
		return errors.New(ErrorUploadFile)

	}
	if err != nil {
		slog.Error("UploadSecure; an unexpected error happened during uploading", "ERROR", err)
		return errors.New(ErrorStrangeUploadFile)
	}
	return errors.New(ErrorStrangeUploadFile)
}
func (sa *NewUploading) UploadFile(data UploadFileIncomingData) error {
	logger := slog.With("UploadFile")
	uploader := manager.NewUploader(sa.S3Info.S3Connect, func(uploader *manager.Uploader) {
		uploader.MaxUploadParts = 1000
		uploader.PartSize = int64(data.Parts * 1024 * 1024)
		uploader.Concurrency = data.Goroutines
	})

	logger.Info("File uploading details", slog.Group("Info", slog.String("File extension", data.FileFormat),
		slog.Int("Parts", data.Parts), slog.Int("Goroutines", data.Goroutines)))
	_, err := uploader.Upload(data.Ctx, &s3.PutObjectInput{
		Bucket:      aws.String(sa.S3Info.Bucket),
		Key:         aws.String(data.FileName),
		ContentType: aws.String(data.FileFormat),
		Body:        data.FileBody.Normal,
	})

	switch {
	case errors.Is(err, context.Canceled):
		logger.Error("The user stopped uploading", "ERROR", err)
		return errors.New(ErrorUploadFile)

	}
	if err != nil {
		logger.Error("an unexpected error", "ERROR", err)
		return errors.New(ErrorStrangeUploadFile)
	}
	return nil
}
