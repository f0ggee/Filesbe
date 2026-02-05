package Handlers

import (
	"Kaban/iternal/InfrastructureLayer"
	Uttiltesss2 "Kaban/iternal/Service/Helpers"
	"Kaban/iternal/Service/Helpers/validator"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func FileUploaderNoEncrypt(r *http.Request) (string, error) {
	slog.Info("Func FileUploaderNoEncrypt starts")

	apps := *InfrastructureLayer.ConnectKeyControl()

	file, sizeAndName, err := r.FormFile("file")
	if err != nil {
		slog.Error("Err from FileUploader 1 ", err)
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

	err = validator.CheckFileSize2(sizeAndName.Size)
	if err != nil {
		slog.Error("Error in file Uploader no encrypt", "Error", err)
		return "", err
	}

	ctx, cancel := Uttiltesss2.Context2(r.Context())
	if cancel == nil {
		return "", errors.New("error in file Uploader no encrypt")
	}
	defer cancel()

	shortNameFile := apps.Key.GenerateShortFileName()

	_, goroutines := Uttiltesss2.FindBest(sizeAndName.Size)

	timeS := time.Now()

	defer func() {
		sa := time.Since(timeS)
		slog.Info("Time of downloading", "Time", sa)
	}()

	cfg, err := Uttiltesss2.S3Helper()
	if err != nil {
		slog.Error("Err cant", err)
		return "", err
	}

	s, err2, done := uploadFile(cfg, goroutines, err, ctx, sizeAndName, file)
	if done {
		return s, err2
	}

	fileIntoBytes, err := json.Marshal(sizeAndName.Filename)
	if err != nil {
		slog.Error("Err in FileUploader no encrypt", "Error", err)
		return "", err
	}
	redisConnect := *InfrastructureLayer.NewSetRedisConnect()

	err = redisConnect.Ras.WriteData(shortNameFile, fileIntoBytes)
	if err != nil {
		return "", err
	}

	slog.Info("File success upload")

	return shortNameFile, nil

}

func uploadFile(cfg *s3.Client, goroutines int, err error, ctx context.Context, sizeAndName *multipart.FileHeader, file multipart.File) (string, error, bool) {
	uploader := manager.NewUploader(cfg, func(uploader *manager.Uploader) {
		uploader.MaxUploadParts = 1000
		uploader.PartSize = 50 * 1024 * 1024
		uploader.Concurrency = goroutines
	})

	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(Bucket),
		Key:    aws.String(sizeAndName.Filename),
		Body:   file,
	})

	switch {
	case errors.Is(err, context.Canceled):
		slog.Info("a user has been cancelled download ")
		return "", errors.New("a user has been cancelled download"), true

	}
	if err != nil {
		slog.Error("Error in uploader", err)
		return "", err, true
	}
	return "", nil, false
}
