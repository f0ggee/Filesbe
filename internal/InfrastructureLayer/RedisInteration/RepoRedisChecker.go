package RedisInteration

import (
	"Kaban/internal/Dto"
	"context"
	"log/slog"
	"time"
)

type RedisCheck struct {
}

func NewRedisCheck() RedisCheck {
	return RedisCheck{}
}

func (d *RedisCheck) ChekIsStartDownloadTest(s string, context context.Context) bool {

	select {
	case <-time.After(1 * time.Second):
		if ax, ds := context.Value("IsFllen").(bool); ax {
			if ds {
				return true
			}
			return false
		}
		return false
	case <-context.Done():
		slog.Error("ChekIsStartDownloadTest;Error context status", "Error", context.Err())
		return true
	}

}

func (d *RedisCheck) ChekIsStartDownload(name string, ctx context.Context) bool {

	isExit := Dto.FileInfoLabels{
		InfoAboutFile:   nil,
		IsStartDownload: false,
	}
	err := redisConn.HGetAll(context.Background(), name).Scan(&isExit)

	if err != nil {
		slog.Error("Can't get the label IsStartDownload", "Error", err)
		return false
	}
	if isExit.IsStartDownload {
		return true
	}

	return false

}
func (d *RedisCheck) CheckFileInfoExists(FileName string, ctx context.Context) bool {

	c, err := redisConn.Exists(ctx, FileName).Result()
	if err != nil {
		slog.Error("CheckExistFileInfo error:", "Error", err)
		return false
	}

	if c > 0 {
		return true
	}
	return false
}
