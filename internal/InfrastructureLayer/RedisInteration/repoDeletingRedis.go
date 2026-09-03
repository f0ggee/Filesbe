package RedisInteration

import (
	"context"
	"errors"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type DeleterRedis struct {
	Re *redis.Client
}

func GetNewDeleterRedis(re *redis.Client) *DeleterRedis {
	return &DeleterRedis{Re: re}
}

func (d *DeleterRedis) DeleterFileInfoTest(s string, context context.Context) error {

	if ax, dsa := context.Value("isFallRedis").(bool); ax != false {

		if dsa == true {
			return errors.New("error by Redis.DeleterFileInfoTest")
		}
		return nil
	}
	return nil
}

func (d *DeleterRedis) DeleteFileInfo(fileInfo string, ctx context.Context) error {

	err := d.Re.Del(ctx, fileInfo).Err()
	if err != nil {
		slog.Error("File info's already been deleted", err)
		return errors.New(ErrorDeleteInfo)
	}
	return nil
}
