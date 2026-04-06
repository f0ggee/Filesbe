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
		slog.Error("Unmarshal err", "Error", err.Error())
		return err, ""
	}

	FileBody, err := sa.S3.S3Download.Download(trueFileName, IncomeContext)
	if err != nil {
		return err, ""
	}
	defer func() {
		slog.Info("Downloading was completed")
		FileBody.Body.Close()

	}()

	w.Header().Set("Content-Type", sa.FileInfo.FileManaging.FindFormatOfFile(trueFileName))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename= %v", trueFileName))
	w.Header().Set("Content-Length", strconv.FormatUint(uint64(*FileBody.ContentLength), 10))

	if _, err = io.Copy(w, FileBody.Body); err != nil {
		slog.Error("Err In file Service Downloader", "err", err)
		return errors.New("connect close"), ""

	}

	slog.Info("Start deleting function")
	err = sa.S3.Deleter.DeleteFileFromS3(trueFileName, IncomeContext)
	if err != nil {
		return err, ""
	}
	slog.Info("Finish deleting function")

	return nil, ""
}
