package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/FileControls"
	s3Repo2 "Kaban/internal/InfrastructureLayer/FileTransferring/s3Repo"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type DownloadNetwork struct {
	W http.ResponseWriter
}

type DownloadFileControl struct {
	Transfer   FileControls.Transferring
	Downloader DomainLevel.MakerDownloader
	Uploader   DomainLevel.MakerUploader
	Deleter    DomainLevel.MakerDeleter
}
type DownloadDelivery struct {
	Reader     DomainLevel.ReadingRedis
	DeleterS3  s3Repo2.DeleterS3
	S3Download s3Repo2.DownloadingS3
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
	g, _ := errgroup.WithContext(IncomeContext)

	trueFileName := ""
	err = json.Unmarshal(fileNameInBytes, &trueFileName)
	if err != nil {
		slog.Error("Application Downloader: error to decode data", "ERROR", err.Error())
		return errors.New(DomainLevel.ErrorParseInfo)
	}

	downloaderObject, err := sa.Downloader.SetName(name).Make(IncomeContext)
	if err != nil {
		return err
	}
	uploaded, err := sa.Uploader.SetName(name).SetSize(sa.Downloader.GetFileSize()).SetAddWriter(sa.W).Make(IncomeContext)
	if err != nil {
		return err
	}
	defer uploaded.CloseSource()

	g.Go(func() error {
		dow, err := downloaderObject.Downloader()
		if err != nil {
			return err
		}
		err = uploaded.Uploader(dow)
		return err
	})

	if err = g.Wait(); err != nil {
		return err
	}

	deletedObj, err := sa.Deleter.SetName(trueFileName).Make(IncomeContext)
	if err != nil {
		return err
	}
	err = deletedObj.Deleter()
	if err != nil {
		return err
	}
	return nil
}
