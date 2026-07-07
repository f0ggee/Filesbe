package RedisInteration

import (
	"Kaban/internal/DomainLevel"

	"github.com/redis/go-redis/v9"
)

func ConnectToRedis() *redis.Client {
	redisConnect := redis.NewClient(&redis.Options{
		Addr:     DomainLevel.RedisHost,
		Username: DomainLevel.RedisServer,
		Password: DomainLevel.RedisPassword,
	})
	return redisConnect

}
