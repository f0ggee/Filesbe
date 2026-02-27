package RedisInteration

import (
	"context"
	"log/slog"
)

func (d *RedisInterationLayer) CheckExistFileInfo(FileName string) bool {

	c, err := d.Re.Exists(context.Background(), FileName).Result()
	if err != nil {
		slog.Error("CheckExistFileInfo error:", err)
		return false
	}

	if c > 0 {
		return true
	}
	return false
}
