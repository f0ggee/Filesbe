package Controller

import (
	"Kaban/iternal/Service/Handlers"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const JsonExample = "application/json"

func getNameFromUrl2(r *http.Request) string {
	vars := mux.Vars(r)

	name := vars["name"]
	return name

}
func DownloadWithNotEncrypt(w http.ResponseWriter, r *http.Request, s *Handlers.HandlerPackCollect) {
	type JsonAnswer struct {
		StatusOperation string   `json:"StatusOperation"`
		Error           []string `json:"Error"`
	}
	if r.Method != http.MethodGet {

		w.Header().Set("Content-Type", JsonExample)
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(JsonAnswer{StatusOperation: "BREAK", Error: []string{"Method don't allow "}}); err != nil {
			slog.Error("Error parse json in answer", err)
			return
		}
		return
	}

	name := getNameFromUrl2(r)

	err, _ := s.DownloadWithNonEncrypt(w, name, r.Context())

	switch {
	case strings.Contains(fmt.Sprint(err), "file was used"):
		slog.Error("File was used and we do redirect to the information page")
		http.Redirect(w, r, "/informationPage", http.StatusFound)
		return

	}
	if err != nil {
		slog.Error("Error downloading file", err)

		return

	}
	return

}
