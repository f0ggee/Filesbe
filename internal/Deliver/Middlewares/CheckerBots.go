package Middlewares

import (
	http2 "Kaban/internal/Deliver/httpController"
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

func CheckBots(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Answer struct {
			StatusOperation string `json:"StatusOperation"`
			Error           string `json:"Error"`
			UrlToRedict     string `json:"UrlRedict"`
		}
		UserAgent := r.Header.Get("User-Agent")
		if strings.Contains(UserAgent, DomainLevel.Bots) {
			w.Header().Set("Content-Type", DomainLevel.Json)
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(&Answer{
				StatusOperation: http2.Break,
				Error:           "Request isn't correct",
				UrlToRedict:     "",
			}); err != nil {

				slog.Error("Error", "ERROR", err.Error())
				return
			}

			return
		}
		next.ServeHTTP(w, r)
	})
}
