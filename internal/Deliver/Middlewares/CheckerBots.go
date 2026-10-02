package Middlewares

import (
	"Kaban/internal/Deliver/httpController"
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
		if strings.Contains(UserAgent, httpController.Bots) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(&Answer{
				StatusOperation: httpController.Break,
				Error:           "Request isn't correct",
				UrlToRedict:     "",
			}); err != nil {
				slog.Error("CheckBots:error to check", "ERROR", err.Error())
				return
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}
