package Helpers

import (
	"context"
	"time"
)

func ContextForDownloading(ctx context.Context) (ctxe context.Context, cancel context.CancelFunc) {

	return context.WithTimeout(ctx, 15*time.Minute)

}

func Context2(IncomingRequest context.Context, n time.Duration) (ctx context.Context, cancel context.CancelFunc) {

	return context.WithTimeout(IncomingRequest, n)

}
