package RedisChecking

import "github.com/redis/go-redis/v9"

type ValidationRedis struct {
	Re *redis.Client
}
