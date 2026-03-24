package Handlers

import (
	"context"
	"encoding/json"
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

func (sa *HandlerPackCollect) DownloadWithNonEncrypt(w http.ResponseWriter, name string, IncomeContext context.Context) (error, string) {

	slog.Info("Func DownloadWithNonEncrypt starts")

	fileNameInBytes, err := sa.RedisControlling.Reader.GetFileInfo(name, IncomeContext)
	if err != nil {
		return err, ""
	}

	trueFileName := ""
	err = json.Unmarshal(fileNameInBytes, &trueFileName)
	if err != nil {
		slog.Error("Unmarshal err", err.Error())
		return err, ""
	}

	downloader := s3.New(sa.S3.S3OldConnect)

	o, err := downloader.GetObjectWithContext(IncomeContext, &s3.GetObjectInput{
		Bucket:      aws.String(Bucket),
		IfNoneMatch: aws.String(""),

		Key: &trueFileName,
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

	w.Header().Set("Content-Type", sa.FileInfo.FileManaging.FindFormatOfFile(trueFileName))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename= %v", trueFileName))
	w.Header().Set("Content-Length", strconv.FormatInt(*o.ContentLength, 10))

	if _, err = io.Copy(w, o.Body); err != nil {
		slog.Error("Err In file Service Downloader", "err", err)
		return errors.New("connect close"), ""

	}

	slog.Info("start delete func in download  ")

	err = sa.S3.Deleter.DeleteFileFromS3(trueFileName, context.Background())
	if err != nil {
		return err, ""
	}
	slog.Info("ends delete func in download  ")

	return nil, ""
}
