package ReadingRedis

import (
	"Kaban/internal/Dto"
	"context"
	"log/slog"
)

func (d *RedisReader) GetFileInfo(fileInfoName string, ctx context.Context) ([]byte, error) {

	StructOfFileInfo := Dto.FileInfoLabels{
		InfoAboutFile: nil,
	}

	err := d.Re.HGetAll(ctx, fileInfoName).Scan(&StructOfFileInfo)
	if err != nil {
		slog.Error("Error in  read data", err)
		return nil, err
	}

	return StructOfFileInfo.InfoAboutFile, nil

}
