package RepoRegisterRepository

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"log/slog"
	"net/http"
)

type NewRegister struct{}

var ErrorController = &NewRegister{}

type RegisterAnswers interface {
	ErrorAnswer(RegisterErrorIncomingData)
	GoodAnswer(RegisterGoodIncomingData)
}
type RegisterErrorIncomingData struct {
	W         http.ResponseWriter
	Error     error
	Operation error
}
type RegisterGoodIncomingData struct {
	W         http.ResponseWriter
	Operation string
	Redirect  string
}

func GetNewRegisterController() *NewRegister {
	return &NewRegister{}
}

func (r NewRegister) ErrorAnswer(data RegisterErrorIncomingData) {
	userErrorResponse := DomainLevel.RegisterAnswer{
		StatusOfOperation: data.Operation.Error(),
		UrlToRedirect:     "",
		Error:             data.Error.Error(),
	}
	data.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	data.W.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(data.W).Encode(userErrorResponse); err != nil {
		slog.Error("ErrorAnswerRegister; error during encoding", "ERROR", err)
		return
	}
	return
}
func (r NewRegister) GoodAnswer(data RegisterGoodIncomingData) {
	d := DomainLevel.RegisterAnswer{
		StatusOfOperation: data.Operation,
		UrlToRedirect:     data.Redirect,
	}
	data.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	data.W.WriteHeader(http.StatusOK)
	err := json.NewEncoder(data.W).Encode(&d)
	if err != nil {
		slog.Error("RegisterGoodAnswer: the error happened in encoding data for a client", "ERROR", err)
		return
	}
	return
}
