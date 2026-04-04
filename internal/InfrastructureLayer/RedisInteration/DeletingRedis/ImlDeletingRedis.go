package DeletingRedis

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

type DeleterRedis struct {
	Re *redis.Client
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
