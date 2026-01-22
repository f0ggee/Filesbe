package Handlers

import (
	Uttiltesss2 "Kaban/iternal/Service/Helpers"
	"context"
	"errors"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func FileUploaderNoEncrypt(w http.ResponseWriter, r *http.Request) (string, error) {
	slog.Info("Func FileUploaderNoEncrypt starts")

	file, sizeAndName, err := r.FormFile("file")
	if err != nil {
		slog.Error("Err from FileUploader 1 ", err)
		http.Error(w, "Error", http.StatusNotFound)
		return "", err
	}

	ctx, cancel, err := CheckFileSize2(r.Context(), sizeAndName)
	if err != nil {
		slog.Error("Error in file Uploader no encrypt", "Error", err)
		return "", err
	}

	//This function cheks a len of name file
	nameFile := Uttiltesss2.CheckLenOfName(sizeAndName.Filename)
	defer cancel()

	_, goroutines := Uttiltesss2.FindBest(sizeAndName.Size)

	defer func() {
		err = file.Close()
		if err != nil {
			slog.Error("Err, cant' close a file", "err", err)
			return
		}
	}()
	timeS := time.Now()

	defer func() {
		sa := time.Since(timeS)
		slog.Info("Time of downloading", "Time", sa)
	}()

	cfg, err := Uttiltesss2.Initialization2()
	if err != nil {
		slog.Error("Err cant", err)
		return "", err
	}

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
	if err != nil {
		slog.Error("Error in uploader", err)
		return "", err
	}

	slog.Info("File success upload")

	return nameFile, nil

}

func CheckFileSize2(incomingRequest context.Context, sizeAndName *multipart.FileHeader) (context.Context, context.CancelFunc, error) {

	sizeFile := sizeAndName.Size
	switch {
	case sizeFile >= 100000000 && sizeFile <= 500000000:
		ctx, c := Uttiltesss2.Context2(incomingRequest, 5*time.Minute)

		return ctx, c, nil

	case sizeFile >= 500000000 && sizeFile < 1000000000:
		ctx, c := Uttiltesss2.Context2(incomingRequest, 5*time.Minute)

		return ctx, c, nil

	case sizeFile > 1000000000:
		return nil, nil, errors.New("file too big")

	default:
		ctx, c := Uttiltesss2.ContextForDownloading(incomingRequest)
		return ctx, c, nil

	}

}
