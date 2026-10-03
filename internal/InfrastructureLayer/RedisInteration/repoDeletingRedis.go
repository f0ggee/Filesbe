package RedisInteration

import (
	"context"
	"errors"
	"log/slog"
)

type RedisDelete struct {
}

func NewRedisDelete() RedisDelete {
	return RedisDelete{}
}

func (d *RedisDelete) DeleterFileInfoTest(s string, context context.Context) error {

	if ax, dsa := context.Value("isFallRedis").(bool); ax != false {

		if dsa == true {
			return errors.New("error by Redis.DeleterFileInfoTest")
		}
		return nil
	}
	return nil
}

func (d *RedisDelete) DeleteFileInfo(fileInfo string, ctx context.Context) error {

	err := redisConn.Del(ctx, fileInfo).Err()
	if err != nil {
		slog.Error("File info's already been deleted", err)
		return errors.New(ErrorDeleteInfo)
	}
	return nil
}
