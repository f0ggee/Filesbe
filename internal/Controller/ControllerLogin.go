package Controller

import (
	"Kaban/internal/Dto"
	"Kaban/internal/Service/Handlers"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
)

var KeyCookie = []byte{}

func init() {

}

func SessionStore() sessions.Store {

	var store1z, err = hex.DecodeString(os.Getenv("KEY1"))
	if err != nil {
		slog.Error("Err decode the key", "Err", err)
		return nil
	}
	Store := sessions.NewCookieStore(store1z)
	return Store

}

func checkJson(r *http.Request) (*Dto.UserLoginData, error) {
	var err error
	var e Dto.UserLoginData
	slog.Info("Key cookie is ", string(KeyCookie))

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Error is closing the body in the controller login", "Error", err)
			return
		}
	}(r.Body)

	return &e, err
}

func Login(w http.ResponseWriter, r *http.Request, realization *Handlers.HandlerPackCollect) {

	type AnswerLogin struct {
		StatusOfOperation string `json:"StatusOperation"`
		UrlToRedict       string `json:"UrlToRedict"`
		ErrorMessage      string `json:"ErrorMessage"`
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Dont' allow", http.StatusUnauthorized)
		slog.Error("Error", "err")
		return
	}
	store := SessionStore()
	Session, err := store.Get(r, TokenName)
	if err != nil {

		slog.Error("cookie don't send 1 ", err)
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(AnswerLogin{
			StatusOfOperation: Break,
			UrlToRedict:       "",
			ErrorMessage:      "User is not unauthorized",
		}); err != nil {
			slog.Error("Err in controller login", "Err", err)
			w.Header().Set(ContentType, Json)
		}
		return
	}

	sa, err := checkJson(r)
	if err != nil {
		return

	}
	err = ValiDateData(sa)
	if err != nil {
		per := AnswerLogin{
			StatusOfOperation: Break,
			ErrorMessage:      "Data has not been validated",
		}
		w.Header().Set("Content-Type", Json)
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(&per); err != nil {
			ControllerErrorLogger.Error("Json in Login can't treated", "Err", err)
			return

		}
		return

	}

	JwtToken, RefreshToken, err := realization.LoginService(*sa, r.Context())
	if err != nil {
		per := AnswerLogin{
			StatusOfOperation: NotStart,
		}
		w.Header().Set("Content-Type", Json)
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(&per)
		if err != nil {
			ControllerErrorLogger.Error("Json in Login can't treated", "Err", err)
			return
		}
		return
	}

	slog.Info("RefreshToken", "RefreshToken", RefreshToken)
	Session.Values[RTCookieName] = RefreshToken
	Session.Values[JwtCookieName] = JwtToken

	Session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int((100 * time.Hour).Hours()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		//Domain:   r.Host,
	}

	if err := Session.Save(r, w); err != nil {
		return

	}

	w.Header().Set("Content-Type", Json)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(AnswerLogin{
		StatusOfOperation: Success,
		UrlToRedict:       "/main",
	}); err != nil {
		ControllerErrorLogger.ErrorContext(r.Context(), "Json in Login can't treated", "Err", err)
		return

	}

}

func ValiDateData(p *Dto.UserLoginData) error {
	validate := validator.New()

	err := validate.Struct(p)
	if err != nil {
		slog.Error("Can't validate because", "Err", err)
		return err

	}
	return nil
}
