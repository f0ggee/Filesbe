package RepoDownloadNoEncrypt

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type jsonAnswerDownloadNoEncrypt struct {
	StatusOperation string `json:"status_operation"`
	Error           string `json:"error"`
	Url             string `json:"url"`
}

type NewRepoDownloadNoEncrypt struct{}

func GetNewNewRepoDownloadNoEncrypt() *NewRepoDownloadNoEncrypt {
	return &NewRepoDownloadNoEncrypt{}
}

type DownloadNoEncryptIncomingData struct {
	W               http.ResponseWriter
	Err             string
	StatusOperation string
}

type Answer interface {
	SetBadAnswer(data DownloadNoEncryptIncomingData)
}

type UrlWork interface {
	GetData(*http.Request) string
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

type SetDownloadFile interface {
	DownloadFile(DownloadNoEncryptData) error
}

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
type NewDownloadFile struct{}

func (n NewDownloadFile) DownloadFile(details DownloadNoEncryptData) error {
	details.W.Header().Set("Content-Type", details.FileFormat)
	details.W.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename= %v", details.TrueFileName))
	details.W.Header().Set("Content-Length", strconv.FormatUint(uint64(details.FileLength), 10))
	if _, err := io.Copy(details.W, details.FileBody); err != nil {
		slog.Error("EncryptDownloadFile; error to download a file", "ERROR", err)
		return errors.New(DomainLevel.ErrorDownloadFile)
	}
	return nil
}
