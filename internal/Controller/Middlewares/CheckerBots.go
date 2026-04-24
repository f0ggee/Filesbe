package Middlewares

import (
	"Kaban/internal/Controller"
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
			UrlToRedict     string `json:"UrlToRedict"`
		}
		UserAgent := r.Header.Get("User-Agent")
		if strings.Contains(UserAgent, Controller.Bots) {
			w.Header().Set("Content-Type", Controller.Json)
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(&Answer{
				StatusOperation: Controller.Break,
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
