package RedisInteration

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func (*RedisInterationLayer) SetIstartDonwload(nameOfFileInfo string) error {
	redisConnect := ConnectToRedis()
	defer func(redisConnect *redis.Client) {
		err := redisConnect.Close()
		if err != nil {
			slog.Error("Error in redis connect", err)
			return
		}
	}(redisConnect)

	err := redisConnect.HSet(context.Background(), nameOfFileInfo, "IsStartDownload", true).Err()
	if err != nil {
		slog.Error("Error set up the labels isStartDownload on true", err)
		return err
	}

	return nil
}
