package Handlers

import (
	"Kaban/Service/Uttiltesss"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"os"
)

func DownloadWithNonEncrypt(w http.ResponseWriter, name string, IncomeContext context.Context) error {

	ctx, cancle := Uttiltesss.Contexte(IncomeContext)
	defer cancle()

	bucket := os.Getenv("BUCKET")

	sees, err := Uttiltesss.Inzelire()
	if err != nil {
		slog.Error("Error in create s3 server", "err:", err)

		return err
	}

	fmt.Println(os.Getenv("FSA"))
	downloader := s3.New(sees)

	o, err := downloader.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket:      aws.String(bucket),
		IfNoneMatch: aws.String(""),
		Key:         &name,
	})
	defer o.Body.Close()

	if err != nil {
		slog.Error("ServiceDownload:", "err", err)
		return err
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename= %v", name))
	w.Header().Set("Content-Length", strconv.FormatInt(*o.ContentLength, 10))

	if _, err = io.Copy(w, o.Body); err != nil {
		slog.Error("Err In file Service Downloader", "err", errors.New("connect closer"))
		return errors.New("connect close")

	}

	return nil
}
