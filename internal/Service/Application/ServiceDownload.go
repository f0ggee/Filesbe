package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/FileControls"
	"context"
	"io"

	"golang.org/x/sync/errgroup"
)

type DownloadApplication interface {
	Download(context.Context, string) error
}
type DownloadNetwork struct {
	W io.Writer
}

type DownloadFileControl struct {
	Transfer   FileControls.Transferring
	Downloader DomainLevel.MakerDownloader
	Uploader   DomainLevel.MakerUploader
	Deleter    DomainLevel.MakerDeleter
	Decoder    DomainLevel.Decoder
}
type DownloadDelivery struct {
	Reader DomainLevel.ReadingRedis
}
type NewDownload struct {
	DownloadDelivery
	DownloadFileControl
	DownloadNetwork
}

func GetNewDownload(downloadWithNonEncryptDelivery DownloadDelivery, downloadWithNonEncryptFileControl DownloadFileControl, downloadNotEncryptNetwork DownloadNetwork) *NewDownload {
	return &NewDownload{DownloadDelivery: downloadWithNonEncryptDelivery, DownloadFileControl: downloadWithNonEncryptFileControl, DownloadNetwork: downloadNotEncryptNetwork}
}
func (sa *NewDownload) Download(IncomeContext context.Context, name string) error {
	fileNameInBytes, err := sa.Reader.GetFileInfo(name, IncomeContext)
	if err != nil {
		return err
	}
	g, _ := errgroup.WithContext(IncomeContext)

	trueFileName := ""
	err = sa.Decoder.Decode(fileNameInBytes, []byte(trueFileName))
	if err != nil {
		return err
	}

	downloaderObject, err := sa.Downloader.SetName(name).Make(IncomeContext)
	if err != nil {
		return err
	}
	uploaded, err := sa.Uploader.SetName(name).SetSize(sa.Downloader.GetFileSize()).SetAdditionalWriter(sa.W).Make(IncomeContext)
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
		deletedObj, err := sa.Deleter.SetName(trueFileName).Make(IncomeContext)
		if err != nil {
			return err
		}
		err = deletedObj.Deleter()
		if err != nil {
			return err
		}
		return err
	})

	if err = g.Wait(); err != nil {
		return err
	}
	return nil
}
