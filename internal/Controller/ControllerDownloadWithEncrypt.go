package Controller

import (
	"Kaban/internal/Service/Handlers"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func getNameFromUrl(r *http.Request) string {
	vars := mux.Vars(r)

	name := vars["name"]
	return name

}

func DownloadWithEncrypt(w http.ResponseWriter, r *http.Request, s *Handlers.HandlerPackCollect) {
	type JsonAnswer struct {
		StatusOperation string   `json:"StatusOperation"`
		Error           []string `json:"Error"`
		Url             string   `json:"Url"`
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Status method don't allow", http.StatusBadRequest)
		return
	}
	name := getNameFromUrl(r)

	err := s.DownloadEncrypt(w, r.Context(), name)
	if err != nil {
		if err := json.NewEncoder(w).Encode(JsonAnswer{StatusOperation: Break, Error: []string{"File was used"}, Url: "/informationPage"}); err != nil {
		}
		return
	}

	return

}
