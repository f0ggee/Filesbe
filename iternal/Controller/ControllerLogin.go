package Controller

import (
	"Kaban/iternal/Dto"
	"Kaban/iternal/Service/Handlers"
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

func Store() sessions.Store {
	var store1z, err = hex.DecodeString(os.Getenv("KEY1"))
	if err != nil {
		slog.Error("Err decode the key", "Err", err)
		return nil
	}

	Store := sessions.NewCookieStore(store1z)
	return Store

}

func checkJson(r *http.Request) (*Dto.User, error) {
	var err error
	var e Dto.User

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

func SaveTokens(s sessions.Session, Jwt string, RFT string, r *http.Request, w http.ResponseWriter) error {

	return nil

}

func Login(w http.ResponseWriter, r *http.Request) {

	type AnswerLogin struct {
		StatusOfOperation string `json:"StatusOperation"`
		UrlToRedict       string `json:"UrlToRedict"`
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method Dont' allow", http.StatusUnauthorized)
		slog.Error("Error", "err")
		return
	}
	store := Store()

	Session, err := store.Get(r, "token6")
	if err != nil {

		slog.Error("cookie don't send 1 ", err)
		http.Error(w, "cookie dont sen", http.StatusUnauthorized)
		return
	}

	oldRftToken, _ := Session.Values["RT"].(string)

	sa, err := checkJson(r)
	if err != nil {
		slog.Error("Err", "error", err)
		return

	}
	err = ValiDateData(sa)
	if err != nil {
		per := AnswerLogin{
			StatusOfOperation: "BREAK",
		}
		w.Header().Set("Content-Type", JsonExample)
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(&per)
		if err != nil {
			slog.Error("Json in Login can't treated", "Err", err)
		}
		return

	}

	JWT, RFT, err := Handlers.LoginService(*sa, r.Context(), oldRftToken)
	if err != nil {
		per := AnswerLogin{
			StatusOfOperation: "NotStart",
		}
		w.Header().Set("Content-Type", JsonExample)
		http.Error(w, "Cant' processed ", http.StatusConflict)

		err = json.NewEncoder(w).Encode(&per)
		if err != nil {
			slog.Error("Json in Login can't treated", "Err", err)
			return
		}
		return
	}

	Session.Values["RT"] = RFT
	Session.Values["JWT"] = JWT

	Session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int((100 * time.Hour).Hours()),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	if err := Session.Save(r, w); err != nil {
		return

	}
	//err = SaveTokens(*Session, JWT, RFT, r, w)
	//if err != nil {
	//	per := AnswerLogin{
	//		StatusOfOperation: "BREAK",
	//	}
	//	w.Header().Set("Content-Type", JsonExample)
	//	http.Error(w, "Cant' processed ", http.StatusConflict)
	//
	//	err = json.NewEncoder(w).Encode(&per)
	//	if err != nil {
	//		slog.Error("Json in Login can't treated", "Err", err)
	//		return
	//	}
	//	return
	//}

	per := AnswerLogin{
		StatusOfOperation: "SUCCESS",
		UrlToRedict:       "/main",
	}
	w.Header().Set("Content-Type", JsonExample)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&per); err != nil {
		slog.Error("Can't json Encode", "Err", err)

		http.Error(w, "sda", http.StatusUnauthorized)
		return

	}

}

func ValiDateData(p *Dto.User) error {
	validate := validator.New()

	err := validate.Struct(p)
	if err != nil {
		slog.Error("Can't validate because", "Err", err)
		return err

	}
	return nil
}
