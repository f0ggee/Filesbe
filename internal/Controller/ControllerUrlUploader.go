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

func BuildUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		slog.Error("Method don't allow")
		ControllerErrorLogger.ErrorContext(r.Context(), "Method don't allow", "Method", r.Method)
		http.Error(w, "Method doesn't allow ", http.StatusUnauthorized)
		return
	}
	type Answer struct {
		StatusOperation string `json:"StatusOperation"`
		UrlToRedict     string `json:"UrlToRedict"`
		ErrorMessage    string `json:"ErrorMessage"`
	}

	nameFile, bols := UrlUploader(r)
	if nameFile == "" {
		slog.Error("UrlUploader name file empty", "Host", r.Host)
		w.Header().Set("Content-Type", Json)
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(Answer{
			StatusOperation: NotStart,
			UrlToRedict:     "nil",
			ErrorMessage:    "Can't handle the URL",
		}); err != nil {
			ControllerErrorLogger.ErrorContext(r.Context(), "Error collecting the url", "Error", err)
			return
		}
		return
	}

	switch {
	case bols == "true":
		if err := json.NewEncoder(w).Encode(Answer{
			StatusOperation: Success,
			UrlToRedict:     DomainName + "d2/" + nameFile,
			ErrorMessage:    "",
		}); err != nil {
			ControllerErrorLogger.ErrorContext(r.Context(), "Can't handle the URL", slog.Group("Url parameters"),
				slog.Any("Url parameters", r.URL.Query()), slog.String("Type of downloading", bols))
			return
		}
		return

	case bols == "false":
		if err := json.NewEncoder(w).Encode(Answer{
			StatusOperation: Success,
			UrlToRedict:     DomainName + nameFile,
			ErrorMessage:    "",
		}); err != nil {
			return
		}
		ControllerErrorLogger.ErrorContext(r.Context(), "Can't handle the URL", slog.Group("Url parameters"),
			slog.Any("Url parameters", r.URL.Query()), slog.String("Type of downloading", bols))
		w.Header().Set("Content-Type", Json)
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(Answer{
			StatusOperation: Break,
			UrlToRedict:     "nil",
			ErrorMessage:    "The url isn't valid",
		}); err != nil {
			slog.Error("Json in Login can't treated", "Err", err)
		}
	}

}
