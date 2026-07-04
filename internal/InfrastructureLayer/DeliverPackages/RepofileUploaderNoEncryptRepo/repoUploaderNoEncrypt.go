package RepofileUploaderNoEncryptRepo

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type NewUploaderNoEncrypt struct{}

func (n NewUploaderNoEncrypt) UrlBuild(router *mux.Router, fileName string) (string, error) {
	url, err := router.Get("fileName").URL("name", fileName, "bool", "false")
	if err != nil {
		slog.Error("UrlBuild; can't get a file name", "ERROR", err)
		return "", errors.New(DomainLevel.ErrorCantGetFileName)
	}
	return url.Path, nil
}

func GetNewNewUploaderNoEncrypt() *NewUploaderNoEncrypt {
	return &NewUploaderNoEncrypt{}
}

type IncomingData struct {
	W               http.ResponseWriter
	Error           string
	StatusOperation string
	UrlToRedirect   string
}
type Answers interface {
	SetGoodAnswer(IncomingData)
	SetBadAnswer(IncomingData)
}

func (n NewUploaderNoEncrypt) SetGoodAnswer(data IncomingData) {
	data.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	data.W.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(data.W).Encode(DomainLevel.AnswerUrlBuilder{
		Url:             data.UrlToRedirect,
		StatusOperation: data.StatusOperation,
	}); err != nil {
		slog.Error("ERROR fileUploaderEncrypt SetGoodAnswer; error to encode json", "ERROR", err)
		return
	}
	return
}

func (n NewUploaderNoEncrypt) SetBadAnswer(data IncomingData) {
	data.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	data.W.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(data.W).Encode(DomainLevel.AnswerUrlBuilder{
		StatusOperation: data.StatusOperation,
		ErrorMessage:    data.Error,
	}); err != nil {
		slog.Error("ERROR fileUploaderEncrypt SetBadAnswer; error to encode json", "ERROR", err)
		return
	}
	return
}
