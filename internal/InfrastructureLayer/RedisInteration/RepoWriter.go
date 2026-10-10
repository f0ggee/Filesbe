package RedisInteration

import (
	"Kaban/internal/DomainLevel"
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

	err := RedisConn.HSet(data.Ctx, data.FileName, DomainLevel.FileInfoLabels{
		InfoAboutFile:   data.Info,
		IsStartDownload: false,
	}).Err()
	if err != nil {
		slog.Error("WriteData;Redis WriteData; error to write data", "ERROR", err)
		return errors.New(ErrorWrite)
	}
	return nil

}

type RedisWriteTest struct {
}

func NewRedisWriteTest() *RedisWriteTest {
	return &RedisWriteTest{}
}

func (r RedisWriteTest) WriteData(data DomainLevel.WriteDataIncomeData) error {
	return nil
}

func (r RedisWriteTest) EnableDownloadingParameter(s string, ctx context.Context) error {
	return nil
}
