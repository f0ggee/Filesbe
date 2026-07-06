package S3Uploader

import (
	"Kaban/internal/DomainLevel"
	"context"
	"errors"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func (sa *Uploading) UploadFile(data DomainLevel.UploadFileIncomingData) error {
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
		return errors.New(DomainLevel.ErrorUploadFile)

	}
	if err != nil {
		logger.Error("an unexpected error", "ERROR", err)
		return errors.New(DomainLevel.ErrorStrangeUploadFile)
	}
	return nil
}
