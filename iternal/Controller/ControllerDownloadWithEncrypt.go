package Controller

import (
	"Kaban/iternal/Service/Handlers"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func getNameFromUrl(r *http.Request) string {
	vars := mux.Vars(r)

	name := vars["name"]
	return name

}

func CheckBots(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Answer struct {
			StatusOperation string `json:"StatusOperation"`
			Error           string `json:"Error"`
			UrlToRedict     string `json:"UrlToRedict"`
		}
		UserAgent := r.Header.Get("User-Agent")
		if strings.Contains(UserAgent, "Bot") {

			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(&Answer{
				StatusOperation: "ErrorHandingRequest",
				Error:           "Request isn't correct",
				UrlToRedict:     "",
			}); err != nil {
				slog.Error("CheckBots %s", err.Error())
				return
			}

			return
		}
		next.ServeHTTP(w, r)
	})
}
func DownloadWithEncrypt(w http.ResponseWriter, r *http.Request, s *Handlers.HandlerPackCollect) {

	if r.Method != http.MethodGet {
		http.Error(w, "Status method don't allow", http.StatusBadRequest)
		return
	}
	name := getNameFromUrl(r)

	err := s.DownloadEncrypt(w, r.Context(), name)
	if err != nil {
		http.Redirect(w, r, "/informationPage", http.StatusFound)
		return
	}

	return

}
