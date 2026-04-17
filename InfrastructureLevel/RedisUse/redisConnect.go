package RedisUse

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func RedisConnect() *redis.Client {
	redisConnect := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_SERVER"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	return redisConnect
}
