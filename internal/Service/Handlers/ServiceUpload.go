package Handlers

import (
	"context"
	"errors"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"golang.org/x/sync/errgroup"
)

func (sa *HandlerPackCollect) FileUploader(r *http.Request) (string, error) {
	slog.Info("Func FileUploader starts")
	g, ctx := errgroup.WithContext(r.Context())

	file, sizeAndName, err := r.FormFile("file")
	if err != nil {
		slog.Error("Err from FileUploader 1 ", err.Error())
		return "", err
	}
	if sizeAndName.Size >= FileMaxSize {
		slog.Info("File too big")

		return "", errors.New("file too big")
	}
	defer func() {
		err = file.Close()
		if err != nil {
			slog.Error("Err, cant' close a file", "err", err)
			return
		}
	}()

	shortNameFile := sa.Crypto.Generate.GenerateShortName()

	Parts, goroutines := sa.FileInfo.FileManaging.FindBestOptions(sizeAndName.Size)

	timeS := time.Now()

	defer func() {
		sa := time.Since(timeS)
		slog.Info("Time of downloading", "Time", sa)
	}()

	g.Go(func() error {
		select {
		case <-ctx.Done():
			slog.Info("Context done")

			return ctx.Err()
		default:
		}
		err2 := sa.uploadFile(Parts, goroutines, r.Context(), sizeAndName, file)
		if err2 != nil {
			return err2
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return "", err
	}

	fileIntoBytes, err := sa.Convert.Converting.JsonConverter(sizeAndName.Filename)
	if err != nil {
		slog.Error("Err in FileUploader no encrypt", "Error", err)
		return "", err
	}

	err = sa.RedisControlling.Writer.WriteData(shortNameFile, fileIntoBytes, r.Context())
	if err != nil {
		return "", err
	}

	slog.Info("File was generated")

	return shortNameFile, nil

}

func (sa *HandlerPackCollect) uploadFile(parts int, goroutines int, ctx context.Context, sizeAndName *multipart.FileHeader, file multipart.File) error {
	uploader := manager.NewUploader(sa.S3.S3Connect, func(uploader *manager.Uploader) {
		uploader.MaxUploadParts = 1000
		uploader.PartSize = int64(parts * 1024 * 1024)
		uploader.Concurrency = goroutines
	})

	FileExtension := sa.FileInfo.FileManaging.FindFormatOfFile(sizeAndName.Filename)

	slog.Info("File extension", FileExtension)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(Bucket),
		Key:         aws.String(sizeAndName.Filename),
		ContentType: aws.String(FileExtension),

		Body: file,
	})

	switch {
	case errors.Is(err, context.Canceled):
		slog.Info("a user has been cancelled download ")
		return errors.New("a user has been cancelled download")

	}
	if err != nil {
		slog.Error("Error in uploader", err)
		return err
	}
	return nil
}
