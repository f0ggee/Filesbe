package Application

import (
	"Kaban/internal/DomainLevel"
	"context"
	"io"

	"golang.org/x/sync/errgroup"
)

type DownloadFileControl struct {
	Downloader DomainLevel.MakerDownloader
	Uploader   DomainLevel.MakerUploader
	Deleter    DomainLevel.MakerDeleter
	Decoder    DomainLevel.Decoder
}

type NewDownload struct {
	DownloadFileControl
	Reader DomainLevel.ReadingRedis
	W      io.Writer
}

func (sa *NewDownload) DownloadApplication(IncomeContext context.Context, name string) error {
	fileNameInBytes, err := sa.Reader.GetFileInfo(name, IncomeContext)
	if err != nil {
		return err
	}
	g, _ := errgroup.WithContext(IncomeContext)
	var trueFileName string
	err = sa.Decoder.Decode([]byte(trueFileName), fileNameInBytes)
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
