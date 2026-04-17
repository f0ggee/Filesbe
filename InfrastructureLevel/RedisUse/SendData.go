package RedisUse

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisUsing struct {
	Connect *redis.Client
}

func (s *RedisUsing) SendData(data []byte, serverName string) error {

	err := s.Connect.Set(context.Background(), serverName, data, 1*time.Minute).Err()
	if err != nil {
		slog.Error("RedisUsing.SendData()", "Error", err)
		return err
	}
	return nil
}
