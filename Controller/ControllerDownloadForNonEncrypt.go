package Controller

import (
	"Kaban/Service/Handlers"
	"encoding/json"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

const JAson = "application/json"

func CookiGetInControllerDownloaderNoEnc(w http.ResponseWriter, r *http.Request) bool {
	store := Store()

	session, err := store.Get(r, "token6")
	if err != nil {
		slog.Error("cookie don't send", err)
		return false
	}

	rt := session.Values["RT"].(string)
	jwtToken := session.Values["JWT"].(string)

	_, _, err, _ = Auth(rt, jwtToken, session)
	if err != nil {
		slog.Error("Error in file upload",
			"Err", err)
		return false
	}

	return true
}
func getNameFromUrl2(r *http.Request) string {
	vars := mux.Vars(r)

	name := vars["name"]
	slog.Info(name)
	return name

}
func DownloadWithNotEncrypt(w http.ResponseWriter, r *http.Request) {
	type JsonAnser struct {
		StatusOperation string   `json:"StatusOperation"`
		Error           []string `json:"Error"`
	}
	if r.Method != http.MethodGet {

		w.Header().Set("Content-Type", JAson)
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(JsonAnser{StatusOperation: "BREAK", Error: []string{"Method don't allow "}}); err != nil {
			slog.Error("Erro parse json in answer", err)
			return
		}
		return
	}

	name := getNameFromUrl2(r)

	ok := CookiGetInControllerDownloaderNoEnc(w, r)
	if !ok {
		return
	}

	err := Handlers.DownloadWithNonEncrypt(w, name, r.Context())
	if err != nil {
		return
	}

}
