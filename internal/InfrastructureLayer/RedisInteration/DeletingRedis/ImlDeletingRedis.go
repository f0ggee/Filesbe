package DeletingRedis

import "github.com/redis/go-redis/v9"

type DeleterRedis struct {
	Re *redis.Client
}
