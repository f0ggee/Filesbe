package RedisInteration

import "github.com/redis/go-redis/v9"

func ConnectToRedis() *redis.Client {

	redisConnect := redis.NewClient(&redis.Options{
		Addr:     "77.95.206.154:6379",
		Username: "server1",
		Password: "wmE9v(m6-aVEA%",
	})

	return redisConnect

}
