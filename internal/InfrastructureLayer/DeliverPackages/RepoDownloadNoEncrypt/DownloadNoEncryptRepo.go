package RepoDownloadNoEncrypt

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type FileDetails struct {
	FileFormat   string
	TrueFileName string
	FileLength   int64
}
type DownloadNoEncryptData struct {
	W http.ResponseWriter
	FileDetails
	FileBody io.ReadCloser
}

type jsonAnswerDownloadNoEncrypt struct {
	StatusOperation string `json:"status_operation"`
	Error           string `json:"error"`
	Url             string `json:"url"`
}
type DownloadNoEncryptIncomingData struct {
	W               http.ResponseWriter
	Err             string
	StatusOperation string
}

type NewRepoDownloadNoEncrypt struct{}

type Answer interface {
	SetBadAnswer(data DownloadNoEncryptIncomingData)
}
type SetDownloadFile interface {
	DownloadFile(DownloadNoEncryptData) error
}

type UrlWork interface {
	GetData(*http.Request) string
}

func GetNewNewRepoDownloadNoEncrypt() *NewRepoDownloadNoEncrypt {
	return &NewRepoDownloadNoEncrypt{}
}

func (n NewRepoDownloadNoEncrypt) GetData(request *http.Request) string {
	vars := mux.Vars(request)

	name := vars["name"]
	return name
}

func (n NewRepoDownloadNoEncrypt) SetBadAnswer(answer DownloadNoEncryptIncomingData) {
	answer.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	answer.W.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(answer.W).Encode(jsonAnswerDownloadNoEncrypt{
		StatusOperation: answer.StatusOperation,
		Error:           answer.Err,
	}); err != nil {
		slog.Error("DownloadNoEncryptSetBadAnswer; error to encode a response", "ERROR", err)
		return
	}
	return
}
