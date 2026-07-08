package RepoDownloadEncryptRepo

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type answerFileDownloadEncrypt struct {
	StatusOperation string `json:"status_operation"`
	Error           string `json:"error"`
}

type NewFileDownloadEncrypt struct{}

func GetNewFileDownloadEncrypt() *NewFileDownloadEncrypt {
	return &NewFileDownloadEncrypt{}
}

type IncomingDataAnswer struct {
	W             http.ResponseWriter
	Error         string
	UrlToRedirect string
}

type AnswersDownloadEncrypt interface {
	SetBadAnswers(IncomingDataAnswer)
}
type UrlWork interface {
	GetDataRequest(r *http.Request) string
}

func (n NewFileDownloadEncrypt) GetDataRequest(r *http.Request) string {
	vars := mux.Vars(r)

	name := vars["name"]
	return name
}
func (n NewFileDownloadEncrypt) SetBadAnswers(answer IncomingDataAnswer) {
	answer.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	answer.W.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(answer.W).Encode(answerFileDownloadEncrypt{
		Error: answer.Error,
	}); err != nil {
		slog.Error("SetBadAnswersDownloadEncrypt; error to encode", "ERROR", err)
		return
	}
	return
}
