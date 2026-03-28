package Controller

import (
	"Kaban/internal/Service/Handlers"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func FileUploaderEncrypt(w http.ResponseWriter, r *http.Request, router *mux.Router, s *Handlers.HandlerPackCollect) {

	type Answer struct {
		StatusOperation string `json:"StatusOperation"`
		Error           string `json:"Error"`
		UrlToRedict     string `json:"UrlToRedict"`
	}
	if r.Method != http.MethodPost {
		slog.Error("Err in controller uploader")
		w.Header().Set(ContentType, Json)
		err := json.NewEncoder(w).Encode(Answer{
			StatusOperation: NotStart,
			Error:           "method don't allow",

			UrlToRedict: "nil",
		})
		if err != nil {
			return
		}

		return
	}

	err := CookieGet(w, r, s)
	if err != nil {
		w.Header().Set(ContentType, Json)
		w.WriteHeader(401)
		if err := json.NewEncoder(w).Encode(Answer{
			StatusOperation: NotStart,
			Error:           fmt.Sprint(err),

			UrlToRedict: "/login",
		}); err != nil {
			return
		}
		return

	}

	filName, err := s.FileUploaderEncrypt(r)
	if err != nil {
		fmt.Println(err)

		w.Header().Set(ContentType, Json)
		w.WriteHeader(400)
		if err := json.NewEncoder(w).Encode(Answer{
			StatusOperation: NotStart,
			Error:           fmt.Sprint(err),
			UrlToRedict:     "",
		}); err != nil {
			slog.Info("Error in encoding json ", "Error", err)
			return
		}

		return
	}

	url, err := router.Get("fileName").URL("name", filName, "bool", "true")
	if err != nil {
		slog.Error("Error can't treat", err)
		return
	}

	w.Header().Set(ContentType, Json)
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(Answer{StatusOperation: Success,
		Error: "",

		UrlToRedict: url.Path}); err != nil {
		slog.Info("Error in FileUploadingControlling", "Error", err)
		return
	}

}

func CookieGet(w http.ResponseWriter, r *http.Request, s *Handlers.HandlerPackCollect) error {
	//store := SessionStore()

	session, err := SessionStore.Get(r, TokenName)

	if err != nil {
		slog.Error("cookie don't send", err)
		http.Error(w, "cookie dont sen", http.StatusUnauthorized)
		return errors.New("cookie don't set")
	}

	rtToken, ok := session.Values[RTCookieName].(string)
	if !ok {
		return errors.New("cookie dont get RT")
	}
	jwts, _ := session.Values[JwtCookieName].(string)
	jwts, err = s.Auth(rtToken, jwts)
	if err != nil {
		ControllerErrorLogger.ErrorContext(r.Context(), "Error generate a cokkie", err)
		return errors.New("can't validate a tokens")
	}
	if jwts != "" {
		session.Values[JwtCookieName] = jwts
	}
	return nil
}
