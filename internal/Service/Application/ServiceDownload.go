package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/FileControls"
	"Kaban/internal/InfrastructureLayer/s3Repo"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

type DownloadNetwork struct {
	W http.ResponseWriter
}

type DownloadFileControl struct {
	Transfer     FileControls.Transferring
	FileManaging FileControls.FileSettings
}
type DownloadDelivery struct {
	Reader     DomainLevel.ReadingRedis
	DeleterS3  s3Repo.DeleterS3
	S3Download s3Repo.DownloadingS3
}
type NewDownload struct {
	DownloadDelivery
	DownloadFileControl
	DownloadNetwork
}

func GetNewDownload(downloadWithNonEncryptDelivery DownloadDelivery, downloadWithNonEncryptFileControl DownloadFileControl, downloadNotEncryptNetwork DownloadNetwork) *NewDownload {
	return &NewDownload{DownloadDelivery: downloadWithNonEncryptDelivery, DownloadFileControl: downloadWithNonEncryptFileControl, DownloadNetwork: downloadNotEncryptNetwork}
}
func (sa *NewDownload) Download(name string, IncomeContext context.Context) error {
	fileNameInBytes, err := sa.Reader.GetFileInfo(name, IncomeContext)
	if err != nil {
		return err
	}

	trueFileName := ""
	err = json.Unmarshal(fileNameInBytes, &trueFileName)
	if err != nil {
		slog.Error("Unmarshal err", "Error", err.Error())
		return errors.New(DomainLevel.ErrorParseInfo)
	}

	FileBody, err := sa.S3Download.GetDownload(trueFileName, IncomeContext)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Download; error wasn't")
			return
		}
	}(FileBody.Body)

	err = sa.Transfer.TransferToClient(FileControls.TransferIncomingData{
		W: sa.W,
		FileDetails: FileControls.FileDetails{
			FileFormat:   sa.FileManaging.FindFormatOfFile(trueFileName),
			TrueFileName: trueFileName,
			FileLength:   *FileBody.ContentLength,
		},
		FileBody: FileBody.Body,
	})
	if err != nil {
		return err
	}

	err = sa.DeleterS3.DeleteFileFromS3(trueFileName, IncomeContext)
	if err != nil {
		return err
	}
	return nil
}
