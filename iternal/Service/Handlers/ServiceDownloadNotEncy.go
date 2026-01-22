package Handlers

import (
	"Kaban/iternal/Dto"
	Uttiltesss2 "Kaban/iternal/Service/Helpers"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func DownloadWithNonEncrypt(w http.ResponseWriter, name string, IncomeContext context.Context) (error, string) {

	slog.Info("Func DownloadWithNonEncrypt starts")
	ctx, cancel := Uttiltesss2.ContextForDownloading(IncomeContext)
	defer cancel()

	nameOfFile := GetsName(name)

	//Create connect to S3
	sees, err := Uttiltesss2.Inzelire()
	if err != nil {
		slog.Error("Error in create s3 server", "err:", err)

		return err, ""
	}

	downloader := s3.New(sees)

	o, err := downloader.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket:      aws.String(Bucket),
		IfNoneMatch: aws.String(""),
		Key:         &nameOfFile,
	})

	switch {
	case strings.Contains(fmt.Sprint(err), "NoSuchKey"):
		return errors.New("file was used"), ""

	case err != nil:
		slog.Error("ServiceDownload:", "err", err)
		return err, ""

	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Error close the body", "Err", err)
			return
		}
	}(o.Body)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename= %v", nameOfFile))
	w.Header().Set("Content-Length", strconv.FormatInt(*o.ContentLength, 10))

	if _, err = io.Copy(w, o.Body); err != nil {
		slog.Error("Err In file Service Downloader", "err", err)
		return errors.New("connect close"), ""

	}

	slog.Info("Func DownloadWithNonEncrypt ends")

	slog.Error("start ")
	DeleteFile(nameOfFile, false)

	return nil, ""
}

// GetsName gets the real name from the map
func GetsName(name string) string {

	names := ""
	if Name, ok := Dto.NamesToConvert[name]; ok {
		names = Name
	}
	slog.Info("name", names)
	return names
}
