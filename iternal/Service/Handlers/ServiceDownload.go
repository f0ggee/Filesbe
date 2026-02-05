package Handlers

import (
	"Kaban/iternal/InfrastructureLayer"
	Uttiltesss2 "Kaban/iternal/Service/Helpers"
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

func DownloadWithNonEncrypt(w http.ResponseWriter, name string, IncomeContext context.Context) (error, string) {

	slog.Info("Func DownloadWithNonEncrypt starts")
	ctx, cancel := Uttiltesss2.ContextForDownloading(IncomeContext)
	defer cancel()

	redisConnect := *InfrastructureLayer.NewSetRedisConnect()

	fileNameInBytes, err := redisConnect.Ras.GetFileInfo(name)
	if err != nil {
		return err, ""
	}

	trueFileName := ""
	err = json.Unmarshal(fileNameInBytes, &trueFileName)
	if err != nil {
		slog.Error("Unmarshal err", err)
		return err, ""
	}
	sees, err := Uttiltesss2.Inzelire()
	if err != nil {
		slog.Error("Error in create s3 server", "err:", err)

		return err, ""
	}

	downloader := s3.New(sees)

	o, err := downloader.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket:      aws.String(Bucket),
		IfNoneMatch: aws.String(""),
		Key:         &trueFileName,
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
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename= %v", trueFileName))
	w.Header().Set("Content-Length", strconv.FormatInt(*o.ContentLength, 10))

	if _, err = io.Copy(w, o.Body); err != nil {
		slog.Error("Err In file Service Downloader", "err", err)
		return errors.New("connect close"), ""

	}

	slog.Info("start delete func in download  ")

	S3Interaction := *InfrastructureLayer.NewConnectToS3()

	err = S3Interaction.Manage.DeleteFileFromS3(name, Bucket)
	if err != nil {
		return err, ""
	}
	slog.Info("ends delete func in download  ")

	return nil, ""
}
