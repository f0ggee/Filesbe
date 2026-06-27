package Dto

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

type UserLoginData struct {
	Email    string `validate:"email,min=2,max=40"`
	Password string `validate:"required,min=6"`
}

func (r *UserLoginData) ParseData(request *http.Request) error {
	if err := json.NewDecoder(request.Body).Decode(&r); err != nil {
		slog.Error("GetIncomingData: the error to parse data", "ERROR", err)
		return errors.New(DomainLevel.ErrorParseInfo)

	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Error is closing the body in the controller register", "Error", err)
			return
		}
	}(request.Body)

	return nil

}
