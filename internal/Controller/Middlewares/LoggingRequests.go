package Middlewares

import (
	"Kaban/internal/Controller"
	"context"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/exp/rand"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		slog.Group("Request's info", "USER-AGENT", r.Header.Get("User-Agent"), "Request's type", r.Method, slog.String("Time", time.Now().Format(time.DateTime)))
		ctx := context.WithValue(r.Context(), Controller.RequestId, rand.Int())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
