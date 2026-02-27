package RedisInteration

import (
	"Kaban/iternal/Dto"
	"context"
	"log/slog"
)

func (d *RedisInterationLayer) ChekIsStartDownload(name string) bool {

	isExit := Dto.FileInfo{
		InfoAboutFile:   nil,
		IsStartDownload: false,
	}
	err := d.Re.HGetAll(context.Background(), name).Scan(&isExit)

	if err != nil {
		slog.Error("Can't get the label IsStartDownload", err)
		return false
	}
	if isExit.IsStartDownload {
		return true
	}

	return false

}
