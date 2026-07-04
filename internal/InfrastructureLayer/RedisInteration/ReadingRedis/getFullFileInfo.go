package ReadingRedis

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"context"
	"errors"
	"log/slog"
)

func (d *RedisReader) GetFileInfo(fileInfoName string, ctx context.Context) ([]byte, error) {

	StructOfFileInfo := Dto.FileInfoLabels{
		InfoAboutFile: nil,
	}

	err := d.Re.HGetAll(ctx, fileInfoName).Scan(&StructOfFileInfo)
	if err != nil {
		slog.Error("Redis GetFileInfo; error happened during getting info about a file", "ERROR", err)
		return nil, errors.New(DomainLevel.ErrorFindFileInfo)
	}

	return StructOfFileInfo.InfoAboutFile, nil

}
