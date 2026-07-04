package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadNoEncrypt"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

type DownloadNotEncryptNetwork struct {
	W http.ResponseWriter
}

type NewDownloadNotEncrypt struct {
	RedisControlling
	S3Controlling
	HandlerFileManagerPack
	FileDownload
	DownloadNotEncryptNetwork
}

func GetNewNewDownloadNotEncrypt(redisControlling RedisControlling, s3Controlling S3Controlling, handlerFileManagerPack HandlerFileManagerPack, fileDownload FileDownload) *NewDownloadNotEncrypt {
	return &NewDownloadNotEncrypt{RedisControlling: redisControlling, S3Controlling: s3Controlling, HandlerFileManagerPack: handlerFileManagerPack, FileDownload: fileDownload}
}

func (sa *NewDownloadNotEncrypt) DownloadWithNonEncrypt(name string, IncomeContext context.Context) error {
	fileNameInBytes, err := sa.RedisControlling.Reader.GetFileInfo(name, IncomeContext)
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
			slog.Error("DownloadWithNonEncrypt; error wasn't")
			return
		}
	}(FileBody.Body)

	err = sa.Download.DownloadFile(RepoDownloadNoEncrypt.DownloadNoEncryptData{
		W: sa.W,
		FileDetails: RepoDownloadNoEncrypt.FileDetails{
			FileFormat:   sa.FileManaging.FindFormatOfFile(trueFileName),
			TrueFileName: trueFileName,
			FileLength:   *FileBody.ContentLength,
		},
		FileBody: FileBody.Body,
	})
	if err != nil {
		return err
	}

	err = sa.S3Controlling.Deleter.DeleteFileFromS3(trueFileName, IncomeContext)
	if err != nil {
		return err
	}
	return nil
}
