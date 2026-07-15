package RedisInteration

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"context"
	"errors"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type Writing struct {
	Re *redis.Client
}

func GetNewWriting(re *redis.Client) *Writing {
	return &Writing{Re: re}
}

func (d *Writing) EnableDownloadingParameter(nameOfFileInfo string, ctx context.Context) error {

	err := d.Re.HSet(ctx, nameOfFileInfo, "IsStartDownload", true).Err()
	if err != nil {
		slog.Error("EnableDownloadingParameter;Error set up the labels isStartDownload on true", "ERROR", err.Error())
		return err
	}

	return nil
}

func (s *Writing) WriteData(data DomainLevel.WriteDataIncomeData) error {

	err := s.Re.HSet(data.Ctx, data.FileName, Dto.FileInfoLabels{
		InfoAboutFile:   data.Info,
		IsStartDownload: false,
	}).Err()
	if err != nil {
		slog.Error("WriteData;Redis WriteData; error to write data", "ERROR", err)
		return errors.New(DomainLevel.ErrorWrite)
	}
	return nil

}
