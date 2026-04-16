package ReadingRedis

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"
)

func (d *RedisReader) GetKey(Ctx context.Context) ([]byte, error) {
	count, sec := 0, 0

	ctx, cancel := context.WithTimeout(Ctx, time.Second*10)
	defer cancel()
	for {
		if count > 20 {
			return nil, errors.New("timeout")
		}
		err := d.Re.Get(ctx, os.Getenv("serverName")).Err()

		if err != nil {
			count, sec = +1, +1
			time.Sleep(time.Duration(sec) * time.Second)
			continue
		}
		var data []byte
		err = d.Re.Get(ctx, os.Getenv("serverName")).Scan(&data)
		if err != nil {
			slog.Error("We got the error when try get the data", "Error", err)
			return nil, errors.New(err.Error())
		}

		err = d.Re.Del(ctx, os.Getenv("serverName")).Err()
		if err != nil {
			slog.Error("We got the error", "Error", err)
			return nil, errors.New(err.Error())
		}

		slog.Group("Data was gotten", "info",
			slog.Int("count", count),
			slog.Int("sec", sec))
		return data, nil

	}
}
