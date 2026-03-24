package RedisChecking

import (
	"Kaban/internal/Dto"
	"context"
	"log/slog"
)

func (d *ValidationRedis) ChekIsStartDownload(name string, ctx context.Context) bool {

	isExit := Dto.FileInfoLabels{
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
