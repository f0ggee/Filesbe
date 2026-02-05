package RedisInteration

import (
	"Kaban/iternal/Dto"
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func (*RedisInterationLayer) GetFileInfo(fileInfoName string) ([]byte, error) {

	redisConnect := ConnectToRedis()
	defer func(redisConnect *redis.Client) {
		err := redisConnect.Close()
		if err != nil {
			slog.Error("Error in redis connect", err)
			return
		}
	}(redisConnect)

	StructOfFileInfo := Dto.FileInfo{
		InfoAboutFile: nil,
	}

	err := redisConnect.HGetAll(context.Background(), fileInfoName).Scan(&StructOfFileInfo)
	if err != nil {
		slog.Error("Error in  read data", err)
		return nil, err
	}

	return StructOfFileInfo.InfoAboutFile, nil

}
