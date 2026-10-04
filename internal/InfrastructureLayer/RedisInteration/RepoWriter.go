package RedisInteration

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"context"
	"errors"
	"log/slog"
)

type RedisWrite struct {
}

func NewRedisWrite() RedisWrite {
	return RedisWrite{}
}

func (d *RedisWrite) EnableDownloadingParameter(nameOfFileInfo string, ctx context.Context) error {

	err := RedisConn.HSet(ctx, nameOfFileInfo, "IsStartDownload", true).Err()
	if err != nil {
		slog.Error("EnableDownloadingParameter;Error set up the labels isStartDownload on true", "ERROR", err.Error())
		return err
	}

	return nil
}

func (s *RedisWrite) WriteData(data DomainLevel.WriteDataIncomeData) error {

	err := RedisConn.HSet(data.Ctx, data.FileName, Dto.FileInfoLabels{
		InfoAboutFile:   data.Info,
		IsStartDownload: false,
	}).Err()
	if err != nil {
		slog.Error("WriteData;Redis WriteData; error to write data", "ERROR", err)
		return errors.New(ErrorWrite)
	}
	return nil

}
