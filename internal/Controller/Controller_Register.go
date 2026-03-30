package Controller

import (
	"Kaban/internal/Dto"
	"Kaban/internal/Service/Handlers"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
)

func checkJsonRegister(r *http.Request) (*Dto.UserDataRegister, error) {

	var err error
	var e Dto.UserDataRegister

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
		http.Error(w, "Method don't allow", http.StatusNotFound)
		return
	}

	type RegisterAnswer struct {
		StatusOfOperation string `json:"StatusOfOperation"`
		UrlToRedict       string `json:"UrlToRedict"`
		Error             string `json:"Error"`
	}

	DataRegister, err := checkJsonRegister(r)
	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(RegisterAnswer{
			StatusOfOperation: NotStart,
			UrlToRedict:       "",
			Error:             "Something went wrong with data",
		}); err != nil {
			slog.Error("Error is closing the body in the controller register", "Error", err)
			return
		}

		slog.Error("Error from Controller_register", "err", err)
		return
	}
	err = ValiDateDataForRegister(DataRegister)
	if err != nil {

		w.Header().Set(ContentType, Json)
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(RegisterAnswer{
			StatusOfOperation: Break,
			UrlToRedict:       "",
			Error:             err.Error(),
		}); err != nil {
			slog.Error("Error sesseion", err, "ID", r.Context().Value(RequestId))
			return
		}
		return
	}

	jwt, rt, err := s.RegisterService(DataRegister, r.Context())

	switch {
	case errors.Is(err, errors.New("person already exist")):
		w.Header().Set(ContentType, Json)
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
		slog.Error("Error session", err)

		w.Header().Set(ContentType, Json)
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(RegisterAnswer{
			StatusOfOperation: Break,
			UrlToRedict:       "",
		}); err != nil {
			slog.Error("Error is  Processing the json register response", "Error", err)
			return
		}
		return
	}

	err = NewSession(w, r, jwt, rt)
	if err != nil {
		slog.Error("Error creating session", "Error", err)
		w.Header().Set(ContentType, Json)
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
	w.Header().Set(ContentType, Json)
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

func NewSession(w http.ResponseWriter, r *http.Request, jwt string, rt string) error {
	store := SessionStore()
	session, err := store.Get(r, "token6")
	if err != nil {
		slog.Error("Error get session", err)
		return err

	}
	session.Values[JwtCookieName] = jwt
	session.Values[RTCookieName] = rt
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int((100 * time.Hour).Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Domain:   r.Host,
	}

	if err := session.Save(r, w); err != nil {
		slog.Error("Error in save cookie", "Err", err)
		return err

	}
	return nil

}
func ValiDateDataForRegister(p *Dto.UserDataRegister) error {
	validating := validator.New()

	err := validating.Struct(p)
	if err != nil {
		slog.Error("Can't validate because", "Err", err)
		errsa := err.(validator.ValidationErrors)
		return errsa

	}
	return nil
}
