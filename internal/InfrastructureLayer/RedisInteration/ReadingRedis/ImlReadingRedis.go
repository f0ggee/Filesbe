package ReadingRedis

import "github.com/redis/go-redis/v9"

type RedisReader struct {
	Re *redis.Client
}
