package Controller

import (
	"Kaban/internal/Service/Handlers"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func getNameFromUrl2(r *http.Request) string {
	vars := mux.Vars(r)

	name := vars["name"]
	return name

}
func DownloadWithNotEncrypt(w http.ResponseWriter, r *http.Request, s *Handlers.HandlerPackCollect) {
	type JsonAnswer struct {
		StatusOperation string   `json:"StatusOperation"`
		Error           []string `json:"Error"`
		Url             string   `json:"Url"`
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", Json)
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(JsonAnswer{StatusOperation: Break, Error: []string{"Method don't allow"}}); err != nil {
			slog.Error("Error parse json in answer", err)
			return
		}
		return
	}

	name := getNameFromUrl2(r)

	err, _ := s.DownloadWithNonEncrypt(w, name, r.Context())

	switch {
	case strings.Contains(fmt.Sprint(err), "file was used"):
		slog.Error("Error sesseion", err, "ID", r.Context().Value(RequestId))
		//w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, "/informationPage", http.StatusFound)
		return

	}
	if err != nil {
		ControllerErrorLogger.ErrorContext(r.Context(), "Error downloading file", err)
		if err := json.NewEncoder(w).Encode(JsonAnswer{StatusOperation: Break, Error: []string{"File was used"}, Url: "/informationPage"}); err != nil {
		}
		return
	}
	return

}
