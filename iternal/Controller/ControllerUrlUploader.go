package Controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func UrlUploader(r *http.Request) (string, string) {
	name := r.URL.Query().Get("name")
	bols := r.URL.Query().Get("bool")

	return name, bols
}

func CUrlUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		slog.Error("Method don't allow")
		http.Error(w, "Method", http.StatusUnauthorized)
	}

	nameFile, bols := UrlUploader(r)

	if nameFile == "" {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	slog.Info("as", bols)
	if bols == "true" {
		if err := json.NewEncoder(w).Encode(map[string]string{
			"Url": "https://filesbes.com/" + "d2/" + nameFile,
		}); err != nil {
			slog.Error("Json can't be treated -", err)
			return
		}
		return
	}

	if bols == "false" {
		if err := json.NewEncoder(w).Encode(map[string]string{
			"Url": "https://filesbes.com/" + "d/" + nameFile,
		}); err != nil {
			slog.Error("Json can't be treated -", err)
			return
		}
		return
	}
	slog.Error("Error don't true/false")
	http.Error(w, "Can't be treated because seetins was changed ", http.StatusBadRequest)

}
