package Controller

import (
	"Kaban/iternal/Service/Handlers"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func CheckDenyList(RT string) bool {

	DenyLists := DenyList
	if _, ok := DenyLists[RT]; !ok {
		return true
	}
	return false
}

func getNameFromUrl(r *http.Request) string {
	vars := mux.Vars(r)

	name := vars["name"]
	return name

}
func DownloadWithEncrypt(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Status method don't allow", http.StatusBadRequest)
		return
	}

	type Answer struct {
		StatusOperation string `json:"StatusOperation"`
		Error           string `json:"Error"`
		UrlToRedict     string `json:"UrlToRedict"`
	}
	name := getNameFromUrl(r)

	err := Handlers.DownloadEncrypt(w, r.Context(), name)

	switch {

	case strings.Contains(fmt.Sprint(err), "file was used"):
		http.Redirect(w, r, "/informationPage", http.StatusFound)
		return

	case err != nil:
		w.Header().Set("Content-Type", JsonExample)
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(Answer{
			StatusOperation: "BREAK",
			Error:           fmt.Sprint(errors.New("can't download file")),
			UrlToRedict:     "",
		}); err != nil {
			slog.Error("Error Encode json", "err", err)
			return
		}
		return

	}

	return

}
