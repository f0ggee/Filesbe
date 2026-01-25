package Controller

import (
	"Kaban/iternal/Service/Handlers"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func FileUploaderEncrypt(w http.ResponseWriter, r *http.Request, router *mux.Router) {

	type Answer struct {
		StatusOperation string `json:"StatusOperation"`
		Error           string `json:"Error"`
		UrlToRedict     string `json:"UrlToRedict"`
	}
	if r.Method != http.MethodPost {
		slog.Error("Err in controller uploader")
		w.Header().Set("Content-Type", JsonExample)
		err := json.NewEncoder(w).Encode(Answer{
			StatusOperation: "NotStart",
			Error:           "method don't allow",

			UrlToRedict: "nil",
		})
		if err != nil {
			return
		}

		return
	}

	err := CookieGet(w, r)
	if err != nil {
		w.Header().Set("Content-Type", JsonExample)
		w.WriteHeader(401)
		err := json.NewEncoder(w).Encode(Answer{
			StatusOperation: "NotStart",
			Error:           fmt.Sprint(err),

			UrlToRedict: "/login",
		})

		if err != nil {
			return
		}
		return

	}

	filName, err := Handlers.FileUploaderEncrypt(w, r)
	if err != nil {
		w.Header().Set("Content-Type", JsonExample)
		w.WriteHeader(400)
		err = json.NewEncoder(w).Encode(Answer{
			StatusOperation: "NotStart",
			Error:           fmt.Sprint(err),

			UrlToRedict: "nil",
		})
		if err != nil {
			slog.Info("Error in controller ", err)
			return
		}

		return
	}

	url, err := router.Get("fileName").URL("name", filName, "bool", "true")
	if err != nil {
		slog.Error("Error can't treat", err)
		return
	}

	w.Header().Set("Content-Type", JsonExample)
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(Answer{
		StatusOperation: "SUCCES",
		Error:           "",

		UrlToRedict: url.Path,
	})

}

func CookieGet(w http.ResponseWriter, r *http.Request) error {
	store := Store()

	session, err := store.Get(r, "token6")

	if err != nil {
		slog.Error("cookie don't send", err)
		http.Error(w, "cookie dont sen", http.StatusUnauthorized)
		return errors.New("cookie don't set")
	}

	rtToken, _ := session.Values["RT"].(string)

	jwts, _ := session.Values["JWT"].(string)
	_, _, err = Handlers.Auth(rtToken, jwts)
	if err != nil {
		slog.Error("Auth error", err)
		return errors.New("can't validate a tokens")
	}

	return nil
}
