package Handlers

import (
	"Kaban/iternal/Dto"
	Uttiltesss2 "Kaban/iternal/Service/Uttiltesss"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func FileUploaderNoEncr(w http.ResponseWriter, r *http.Request) (string, error) {

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
	namefile := CheckLenOfName(sizeAndName.Filename)
	defer cancel()

	size, groutine := FindBest(sizeAndName.Size)

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
		fmt.Println(sa)
	}()

	cfg, err := Uttiltesss2.Inzelire2()
	if err != nil {
		slog.Error("Err cant", err)
		return "", err
	}

	bucket := "0c8f1ea9-b07f5996-b392-4227-961b-14d2a71a53dc"
	uploader := manager.NewUploader(cfg, func(uploader *manager.Uploader) {
		uploader.MaxUploadParts = 1000
		uploader.PartSize = int64(size) * 1024 * 1024
		uploader.Concurrency = groutine
	})

	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(sizeAndName.Filename),
		Body:   file,
	})
	if err != nil {
		slog.Error("Error in uploader", err)
		return "", err
	}

	slog.Info("File success upload :)")

	return namefile, nil

}

func CheckLenOfName(sizeAndName string) string {
	nameOfFile := sizeAndName
	if len(sizeAndName) > 5 {
		NewString := rand.Text()
		Dto.NamesToConvert[NewString[:3]] = sizeAndName
		nameOfFile = NewString[:3]

	}
	return nameOfFile
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
		ctx, c := Uttiltesss2.Contexte(incomingRequest)
		return ctx, c, nil

	}

}
