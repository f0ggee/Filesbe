package RedisInteration

import (
	"context"
	"log/slog"
)

func (*RedisInterationLayer) CheckExistFileInfo(FileName string) bool {
	redisConnect := ConnectToRedis()
	defer redisConnect.Close()

	c, err := redisConnect.Exists(context.Background(), FileName).Result()
	if err != nil {
		slog.Error("CheckExistFileInfo error:", err)
		return false
	}

	if c > 0 {
		return true
	}
	return false
}
