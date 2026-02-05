package RedisInteration

import (
	"context"

	"golang.org/x/exp/slog"
)

func (*RedisInterationLayer) DeleteFileInfo(nameFileinfo string) error {

	redisConenct := ConnectToRedis()
	defer redisConenct.Close()
	err := redisConenct.Del(context.Background(), nameFileinfo).Err()
	if err != nil {
		slog.Error("Error in func deleteFileInfo in Redis", err)
		return err
	}
	return nil
}
