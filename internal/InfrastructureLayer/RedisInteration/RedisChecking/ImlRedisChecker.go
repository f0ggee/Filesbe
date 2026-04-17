package RedisChecking

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type ValidationRedis struct {
	Re *redis.Client
}

func (d *ValidationRedis) ChekIsStartDownloadTest(s string, context context.Context) bool {

	select {
	case <-time.After(1 * time.Second):
		slog.Info("Validator starting download test for " + s)
		if ax, ds := context.Value("IsFllen").(bool); ax {
			if ds {
				return true
			}
			return false
		}
		return false
	case <-context.Done():
		slog.Error("Erro context status", context.Err())
		return true
	}

}
