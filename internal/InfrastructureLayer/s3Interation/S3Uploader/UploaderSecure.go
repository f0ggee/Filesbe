package S3Uploader

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

func (sa *Uploading) UploadFileEncrypt(BesParts int, goroutine int, ctx context.Context, shortFileName string, ContentType string, reader *io.PipeReader) error {
	slog.Group("File uploading details",
		slog.String("FileExtension", ContentType),
		slog.String("Parts", fmt.Sprint(BesParts)),
		slog.String("Goroutines", fmt.Sprint(goroutine)),
	)
	uploader := manager.NewUploader(sa.S3Connect, func(uploader *manager.Uploader) {

		uploader.MaxUploadParts = 200
		uploader.PartSize = int64(BesParts) * 1024 * 1024
		uploader.Concurrency = goroutine
		uploader.BufferProvider = manager.NewBufferedReadSeekerWriteToPool(BesParts)
	})

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(sa.Bucket),
		Key:         aws.String(shortFileName),
		ContentType: aws.String(ContentType),
		Body:        reader,
	})
	if err == nil {
		return nil
	}

	var ns *types.NoSuchKey

	switch {

	case errors.As(err, &ns):

		slog.Error("file was used", "Error", err.Error())
		return errors.New("file was used")

	case errors.Is(err, context.Canceled):
		slog.Error("file downloading was cancelled")
		return errors.New("file was cancelled")

	case errors.Is(err, context.DeadlineExceeded):
		slog.Error("Time was exceeded")
		return errors.New("time was exceeded")

	}
	if err != nil {
		slog.Error("Error in file writing", "Error", err)
		return err
	}
	slog.Error("file upload failed", "Error", err)
	return errors.New("file upload failed")
}
