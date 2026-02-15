package RedisInteration

import (
	"Kaban/iternal/Dto"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

func (*RedisInterationLayer) GetKey() ([]byte, []byte, []byte, error) {
	redisConnect := ConnectToRedis()
	defer func(redisConnect *redis.Client) {
		err := redisConnect.Close()
		if err != nil {
			slog.Error("Error closing redis connection", "error", err.Error())
			return
		}
	}(redisConnect)
	count := 0

	for {

		if count > 20 {
			return nil, nil, nil, errors.New("timeout")
		}
		err := redisConnect.HGetAll(context.Background(), "TestServer2").Err()

		if err != nil {
			slog.Error("We got the error", "Error", err)
			count++
			time.Sleep(4 * time.Second)
			continue

		}

		zs := Dto.RedisFileStructFromMasterServer{
			AesKey:    nil,
			PlainText: nil,
			Signature: nil,
		}
		err = redisConnect.HGetAll(context.Background(), "TestServer2").Scan(&zs)
		if err != nil {
			slog.Error("We got the error when try get the data", "Error", err)
			return nil, nil, nil, errors.New(err.Error())
		}

		// TODO - When i start the project, must remove them
		//err = redisConnect.HDel(context.Background(), "server1", "key").Err()
		//if err != nil {
		//	slog.Error("Error delete key", "Error", err)
		//	return nil
		//}

		slog.Info("We got the key", "key", true)
		return zs.AesKey, zs.PlainText, zs.Signature, nil

	}
}
