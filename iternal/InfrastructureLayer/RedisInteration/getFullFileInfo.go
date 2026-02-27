package RedisInteration

import (
	"Kaban/iternal/Dto"
	"context"
	"log/slog"
)

func (d *RedisInterationLayer) GetFileInfo(fileInfoName string) ([]byte, error) {

	StructOfFileInfo := Dto.FileInfo{
		InfoAboutFile: nil,
	}

	err := d.Re.HGetAll(context.Background(), fileInfoName).Scan(&StructOfFileInfo)
	if err != nil {
		slog.Error("Error in  read data", err)
		return nil, err
	}

	return StructOfFileInfo.InfoAboutFile, nil

}
