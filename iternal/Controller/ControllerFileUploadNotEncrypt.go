package Controller

import (
	"Kaban/iternal/Service/Handlers"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func FileUploaderNoEncrypt(w http.ResponseWriter, r *http.Request, router *mux.Router) {
	if r.Method != http.MethodPost {
		http.Error(w, "err", http.StatusUnauthorized)
		slog.Error("Err in Cottroler Uploader")
		return
	}
	type Answer struct {
		StatusOperation string `json:"StatusOperation"`
		UrlToRedict     string `json:"UrlToRedict"`
		Error           string `json:"Error"`
	}

	err := CookiGet2(w, r)
	if err != nil {
		w.Header().Set("Content-Type", JsonExample)
		w.WriteHeader(http.StatusUnauthorized)
		if err = json.NewEncoder(w).Encode(Answer{
			StatusOperation: "NotStart",
			UrlToRedict:     "/login",
		}); err != nil {
			slog.Error("Err in json encode", err)
			return
		}
		return
	}

	filName, err := Handlers.FileUploaderNoEncrypt(r)
	if err != nil {

		w.Header().Set("Content-Type", JsonExample)
		w.WriteHeader(http.StatusUnauthorized)
		if err = json.NewEncoder(w).Encode(Answer{
			StatusOperation: "BREAK",
			UrlToRedict:     "",
			Error:           err.Error(),
		}); err != nil {
			slog.Error("Err in json encode", err)
			return
		}
		return
	}

	url, err := router.Get("fileName").URL("name", filName, "bool", "false")
	if err != nil {
		slog.Error("Error can't treate", err)

		w.Header().Set("Content-Type", JsonExample)
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(Answer{
			StatusOperation: "BREAK",
			UrlToRedict:     "",
		}); err != nil {
			slog.Error("Err in json encode", err)
			return
		}
		return
	}

	w.Header().Set("Content-Type", JsonExample)
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(Answer{
		StatusOperation: "SUCCESS",
		UrlToRedict:     url.Path,
	}); err != nil {
		slog.Error("Err in json encode", err)
		return
	}

}

func CookiGet2(w http.ResponseWriter, r *http.Request) error {
	store := Store()

	session, err := store.Get(r, "token6")
	if err != nil {
		slog.Error("cookie don't send", err)
		http.Error(w, "cookie dont sen", http.StatusUnauthorized)
		return err
	}

	if session.Options.MaxAge == 0 {
		slog.Error("Cookie time expired")
		return errors.New("Cookie time expired")
	}

	rtToken, _ := session.Values["RT"].(string)

	jwts, _ := session.Values["JWT"].(string)
	_, _, err = Handlers.Auth(rtToken, jwts)
	if err != nil {
		slog.Error("cookie doesn't set up ",
			"Err", err)
		return err
	}

	return nil
}
