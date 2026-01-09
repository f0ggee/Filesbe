package Controller

import (
	"Kaban/Dto"
	"Kaban/Service/Handlers"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"log/slog"
	"net/http"
	"time"
)

const ContentType = "Content-Type"

func checkJsonRegister(r *http.Request) (*Dto.Handler_Registerr, error) {

	var err error
	var e Dto.Handler_Registerr

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return nil, err

	}

	defer r.Body.Close()

	return &e, err
}

func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		slog.Error("Error from Controller_register, method don't allow ", "err")
		http.Error(w, "M don't allow", http.StatusNotFound)
		return
	}

	type RegisterAnswer struct {
		StatusOfOperation string `json:"StatusOfOperation"`
		UrlToRedict       string `json:"UrlToRedict"`
	}

	DataRegister, err := checkJsonRegister(r)
	if err != nil {
		slog.Error("Error sesseion", err)

		slog.Error("Error from Controller_register", "err", err)
		return
	}
	err = ValiDateDataForRegister(DataRegister)
	if err != nil {
		slog.Error("Error sesseion", err)

		http.Error(w, "Error treate", http.StatusBadRequest)
		return
	}

	jwt, rt, err := Handlers.RegisterService(DataRegister)
	if err != nil {
		slog.Error("Error sesseion", err)

		w.Header().Set(ContentType, JAson)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RegisterAnswer{
			StatusOfOperation: "BREAK",
			UrlToRedict:       "",
		})
		return
	}

	err = SetSession(w, r, err, jwt, rt)
	if err != nil {
		slog.Error("Error sesseion", err)
		w.Header().Set(ContentType, JAson)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RegisterAnswer{
			StatusOfOperation: "BREAK",
			UrlToRedict:       "",
		})
		return
	}
	w.Header().Set(ContentType, JAson)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RegisterAnswer{
		StatusOfOperation: "SUCCESS",
		UrlToRedict:       "/main",
	})

}

func SetSession(w http.ResponseWriter, r *http.Request, err error, jwt string, rt string) error {
	store := Store()
	session, err := store.Get(r, "token6")
	if err != nil {
		slog.Error("Error get session", err)
		return err

	}
	session.Values["JWT"] = jwt
	session.Values["RT"] = rt
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int((100 * time.Hour).Seconds()),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	if err := session.Save(r, w); err != nil {
		slog.Error("Error in save cookie", "Err", err)
		return err

	}
	return nil

}
func ValiDateDataForRegister(p *Dto.Handler_Registerr) error {
	validater := validator.New()

	err := validater.Struct(p)
	if err != nil {
		slog.Error("Can't validate because", "Err", err)
		return err

	}
	return nil
}
