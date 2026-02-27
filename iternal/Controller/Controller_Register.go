package Controller

import (
	"Kaban/iternal/Dto"
	"Kaban/iternal/Service/Handlers"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
)

const ContentType = "Content-Type"

func checkJsonRegister(r *http.Request) (*Dto.HandlerRegister, error) {

	var err error
	var e Dto.HandlerRegister

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return nil, err

	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Error is closing the body in the controller register", "Error", err)
			return
		}
	}(r.Body)

	return &e, err
}

func Register(w http.ResponseWriter, r *http.Request, s *Handlers.HandlerPackCollect) {

	if r.Method != http.MethodPost {
		slog.Error("Error from Controller_register, method don't allow ", "err")
		http.Error(w, "M don't allow", http.StatusNotFound)
		return
	}

	type RegisterAnswer struct {
		StatusOfOperation string `json:"StatusOfOperation"`
		UrlToRedict       string `json:"UrlToRedict"`
		Error             string `json:"Error"`
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

	jwt, rt, err := s.RegisterService(DataRegister)

	switch {
	case errors.Is(err, errors.New("person already exist")):
		w.Header().Set(ContentType, JsonExample)
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(RegisterAnswer{
			StatusOfOperation: "BREAK",
			UrlToRedict:       "",
			Error:             err.Error(),
		})
		if err != nil {
			slog.Error("Error is  Processing the json register response", "Error", err)
			return
		}
		return
	}
	if err != nil {
		slog.Error("Error sesseion", err)

		w.Header().Set(ContentType, JsonExample)
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(RegisterAnswer{
			StatusOfOperation: "BREAK",
			UrlToRedict:       "",
		})
		if err != nil {
			slog.Error("Error is  Processing the json register response", "Error", err)
			return
		}
		return
	}

	err = SetSession(w, r, err, jwt, rt)
	if err != nil {
		slog.Error("Error sesseion", err)
		w.Header().Set(ContentType, JsonExample)
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(RegisterAnswer{
			StatusOfOperation: "BREAK",
			UrlToRedict:       "",
		})
		if err != nil {
			slog.Error("Error is processing the json register response", "Error", err)
			return
		}
		return
	}
	w.Header().Set(ContentType, JsonExample)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(RegisterAnswer{
		StatusOfOperation: "SUCCESS",
		UrlToRedict:       "/main",
	})
	if err != nil {
		slog.Error("Error is processing the json", "Error", err)
		return
	}

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
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	if err := session.Save(r, w); err != nil {
		slog.Error("Error in save cookie", "Err", err)
		return err

	}
	return nil

}
func ValiDateDataForRegister(p *Dto.HandlerRegister) error {
	validater := validator.New()

	err := validater.Struct(p)
	if err != nil {
		slog.Error("Can't validate because", "Err", err)
		return err

	}
	return nil
}
