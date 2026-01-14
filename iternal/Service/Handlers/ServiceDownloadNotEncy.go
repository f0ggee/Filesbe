package Handlers

import (
	"Kaban/iternal/Dto"
	Uttiltesss2 "Kaban/iternal/Service/Uttiltesss"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"os"
)

func DownloadWithNonEncrypt(w http.ResponseWriter, name string, IncomeContext context.Context) error {

	ctx, cancle := Uttiltesss2.Contexte(IncomeContext)
	defer cancle()

	nameOfFile := name
	nameOfFile = GetsName(name, nameOfFile)

	bucket := os.Getenv("BUCKET")

	sees, err := Uttiltesss2.Inzelire()
	if err != nil {
		slog.Error("Error in create s3 server", "err:", err)

		return err
	}

	downloader := s3.New(sees)

	o, err := downloader.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket:      aws.String(bucket),
		IfNoneMatch: aws.String(""),
		Key:         &nameOfFile,
	})
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Error close the body", "Err", err)
			return
		}
	}(o.Body)

	if err != nil {
		slog.Error("ServiceDownload:", "err", err)
		return err
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename= %v", nameOfFile))
	w.Header().Set("Content-Length", strconv.FormatInt(*o.ContentLength, 10))

	if _, err = io.Copy(w, o.Body); err != nil {
		slog.Error("Err In file Service Downloader", "err", errors.New("connect closer"))
		return errors.New("connect close")

	}

	return nil
}

func GetsName(name string, nameOfFile string) string {
	if Name, ok := Dto.NamesToConvert[name]; ok {
		nameOfFile = Name
		delete(Dto.NamesToConvert, name)
	}
	return nameOfFile
}
