package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Service/Application"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func getNameFromUrl(r *http.Request) string {
	vars := mux.Vars(r)

	name := vars["name"]
	return name

}

func DownloadWithEncrypt(w http.ResponseWriter, r *http.Request, s *Application.HandlerPackCollect) {
	type JsonAnswer struct {
		StatusOperation string   `json:"StatusOperation"`
		Error           []string `json:"Error"`
		Url             string   `json:"Url"`
	}
	if r.Method != http.MethodGet {
		//TODO add handling the error
	}
	name := getNameFromUrl(r)

	err := s.DownloadEncrypt(w, r.Context(), name)

	//TODO add handling the error

	return

}
