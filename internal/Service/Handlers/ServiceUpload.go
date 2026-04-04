package Handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

func (sa *HandlerPackCollect) FileUploader(r *http.Request) (string, error) {
	slog.Info("Func FileUploader starts")
	g, ctx := errgroup.WithContext(r.Context())

	file, sizeAndName, err := r.FormFile("file")
	if err != nil {
		slog.Error("Err from FileUploader 1 ", "Error", err.Error())
		return "", err
	}
	if sizeAndName.Size >= FileMaxSize {
		slog.Info("File too big")

		return "", errors.New("file too big")
	}
	defer func() {
		err = file.Close()
		if err != nil {
			slog.Error("Err, cant' close a file", "err", err)
			return
		}
	}()

	shortNameFile := sa.Crypto.Generate.GenerateShortName()

	Parts, goroutines := sa.FileInfo.FileManaging.FindBestOptions(sizeAndName.Size)

	timeS := time.Now()

	defer func() {
		sa := time.Since(timeS)
		slog.Info("Time of downloading", "Time", sa)
	}()

	fileFormat := sa.FileInfo.FileManaging.FindFormatOfFile(sizeAndName.Filename)
	g.Go(func() error {

		err2 := sa.S3.Uploader.UploadFile(Parts, goroutines, ctx, fileFormat, sizeAndName.Filename, file)
		if err2 != nil {
			return err2
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return "", err
	}

	fileIntoBytes, err := sa.Convert.Converting.JsonConverter(sizeAndName.Filename)
	if err != nil {
		slog.Error("Err in FileUploader no encrypt", "Error", err)
		return "", err
	}

	err = sa.RedisControlling.Writer.WriteData(shortNameFile, fileIntoBytes, r.Context())
	if err != nil {
		return "", err
	}

	slog.Info("File was generated")

	return shortNameFile, nil

}
