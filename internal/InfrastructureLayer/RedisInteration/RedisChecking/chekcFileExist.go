package RedisChecking

import (
	"context"
	"log/slog"
)

func (d *ValidationRedis) CheckFileInfoExists(FileName string, ctx context.Context) bool {

	c, err := d.Re.Exists(ctx, FileName).Result()
	if err != nil {
		slog.Error("CheckExistFileInfo error:", "Error", err)
		return false
	}

	if c > 0 {
		return true
	}
	return false
}
