package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/DeliverHandlers/SessionHandle"
	"Kaban/internal/Service/Application"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func checkJson(r *http.Request) (*Dto.UserLoginData, error) {
	var err error
	var e Dto.UserLoginData

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

func Login(w http.ResponseWriter, r *http.Request, realization *Application.HandlerPackCollect) {

	type AnswerLogin struct {
		StatusOfOperation string `json:"StatusOperation"`
		UrlToRedict       string `json:"UrlRedict"`
		ErrorMessage      string `json:"ErrorMessage"`
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Dont' allow", http.StatusUnauthorized)
		slog.Error("Method Dont' allow", "Method", http.StatusUnauthorized)
		return
	}

	sa, err := checkJson(r)
	if err != nil {
		return

	}
	err = validateData(sa)
	if err != nil {
		per := AnswerLogin{
			StatusOfOperation: DomainLevel.Break,
			ErrorMessage:      "Data has not been validated",
		}
		w.Header().Set("Content-Type", DomainLevel.Json)
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
			StatusOfOperation: DomainLevel.NotStart,
		}
		w.Header().Set("Content-Type", DomainLevel.Json)
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(&per)
		if err != nil {
			ControllerErrorLogger.Error("Json in Login can't treated", "Err", err)
			return
		}
		return
	}

	ReturnedData := SessionHandle.SessionControl.SetNewSession(DomainLevel.IncomingSessionData{
		Writer:  w,
		Request: r,
		Jwt:     JwtToken,
		Rt:      RefreshToken,
	})
	if ReturnedData == nil || ReturnedData.Error != nil {
		//TODO add handling the error
	}
	w.Header().Set("Content-Type", DomainLevel.Json)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(AnswerLogin{
		StatusOfOperation: DomainLevel.Success,
		UrlToRedict:       "/main",
	}); err != nil {
		ControllerErrorLogger.ErrorContext(r.Context(), "Json in Login can't treated", "Err", err)
		return

	}

}

func validateData(p *Dto.UserLoginData) error {
	validate := validator.New()

	err := validate.Struct(p)
	if err != nil {
		slog.Error("validate: can't validate because", "Err", err)
		return err

	}
	return nil
}
