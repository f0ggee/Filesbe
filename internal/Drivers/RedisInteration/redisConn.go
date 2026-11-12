package RedisInteration

import (
	"Kaban/internal/DomainLevel"

	"github.com/redis/go-redis/v9"
)

type redisHub struct {
	RedisConn *redis.Client
}

var hub redisHub

func SetHub(client *redis.Client) {
	hub = redisHub{
		RedisConn: client,
	}
}
func GetRedisConn() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     DomainLevel.RedisHost,
		Username: DomainLevel.RedisServer,
		Password: DomainLevel.RedisPassword,
	})
}
