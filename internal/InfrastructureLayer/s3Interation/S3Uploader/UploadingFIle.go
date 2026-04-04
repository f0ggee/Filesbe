package S3Uploader

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func (sa *Uploading) UploadFile(parts int, goroutines int, ctx context.Context, fileFormat string, fileName string, file multipart.File) error {
	uploader := manager.NewUploader(sa.S3Connect, func(uploader *manager.Uploader) {
		uploader.MaxUploadParts = 1000
		uploader.PartSize = int64(parts * 1024 * 1024)
		uploader.Concurrency = goroutines
	})

	slog.Group("File uploading details",
		slog.String("FileExtension", fileFormat),
		slog.String("Parts", fmt.Sprint(parts)),
		slog.String("Goroutines", fmt.Sprint(goroutines)),
		slog.String("Size", fmt.Sprint()),
	)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(sa.Bucket),
		Key:         aws.String(fileName),
		ContentType: aws.String(fileFormat),
		Body:        file,
	})

	switch {
	case errors.Is(err, context.Canceled):
		slog.Info("a user has been cancelled download", "Error", err)
		return errors.New("a user has been cancelled download")

	}
	if err != nil {
		slog.Error("Error in uploader", "Error", err)
		return err
	}
	return nil
}
