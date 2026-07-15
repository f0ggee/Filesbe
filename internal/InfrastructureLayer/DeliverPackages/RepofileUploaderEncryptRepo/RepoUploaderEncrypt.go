package RepofileUploaderEncryptRepo

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type SetNewUploadingRepo struct{}

func GetNewSetNewUploadingRepo() *SetNewUploadingRepo {
	return &SetNewUploadingRepo{}
}

type IncomingDataAnswer struct {
	W               http.ResponseWriter
	UrlToRedirect   string
	Error           string
	StatusOperation string
}

type answerUploadEncrypt struct {
	StatusOperation string `json:"status_operation"`
	Error           string `json:"error"`
	UrlToRedirect   string `json:"url_to_redirect"`
}

type AnswersUploadEncrypt interface {
	SetGoodAnswers(IncomingDataAnswer)
	SetBadAnswers(IncomingDataAnswer)
}
type UrlUploadEncrypt interface {
	UrlBuilder(r *mux.Router, fileName string) (string, error)
}

func (s SetNewUploadingRepo) UrlBuilder(r *mux.Router, fileName string) (string, error) {
	url, err := r.Get("fileName").URL("name", fileName, "bool", "true")
	if err != nil {
		slog.Error("UrlBuilderUploadEncrypt; error to get a file name from the url", "ERROR", err)
		return "", errors.New(DomainLevel.ErrorCantGetFileName)
	}
	return url.Path, nil
}

func (s SetNewUploadingRepo) SetGoodAnswers(answer IncomingDataAnswer) {

	answer.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	answer.W.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(answer.W).Encode(answerUploadEncrypt{
		StatusOperation: answer.StatusOperation,
		UrlToRedirect:   answer.UrlToRedirect,
	}); err != nil {
		slog.Error("SetGoodAnswersUploadEncrypt; error to encode", "ERROR", err)
		return
	}
	return
}

func (s SetNewUploadingRepo) SetBadAnswers(answer IncomingDataAnswer) {
	answer.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	answer.W.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(answer.W).Encode(answerUploadEncrypt{
		StatusOperation: answer.StatusOperation,
		Error:           answer.Error,
	}); err != nil {
		slog.Error("SetGoodAnswersUploadEncrypt; error to encode", "ERROR", err)
		return
	}
	return
}
