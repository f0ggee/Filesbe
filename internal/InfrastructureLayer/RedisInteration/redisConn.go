package RedisInteration

import (
	"Kaban/internal/DomainLevel"

	"github.com/redis/go-redis/v9"
)

var RedisConn *redis.Client

func init() {
	redisConnect := redis.NewClient(&redis.Options{
		Addr:     DomainLevel.RedisHost,
		Username: DomainLevel.RedisServer,
		Password: DomainLevel.RedisPassword,
	})
	RedisConn = redisConnect
}
