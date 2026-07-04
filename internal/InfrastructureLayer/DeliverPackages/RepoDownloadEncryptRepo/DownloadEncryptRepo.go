package RepoDownloadEncryptRepo

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

type answerFileDownloadEncrypt struct {
	StatusOperation string `json:"status_operation"`
	Error           string `json:"error"`
}

type NewFileDownloadEncrypt struct{}

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

type NewEncryptDownloadFile struct{}
type FileDetails struct {
	FileFormat   string
	TrueFileName string
	FileLength   int64
}
type DownloadEncryptIncomingData struct {
	W http.ResponseWriter
	FileDetails
	FileBody io.ReadCloser
}
type SetDownloadFile interface {
	EncryptDownloadFile(DownloadEncryptIncomingData) error
}

func (n NewEncryptDownloadFile) EncryptDownloadFile(details DownloadEncryptIncomingData) error {
	details.W.Header().Set("Content-Type", details.FileFormat)
	details.W.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename= %v", details.TrueFileName))
	details.W.Header().Set("Content-Length", strconv.FormatUint(uint64(details.FileLength), 10))
	if _, err := io.Copy(details.W, details.FileBody); err != nil {
		slog.Error("EncryptDownloadFile; error to download a file", "ERROR", err)
		return errors.New(DomainLevel.ErrorDownloadFile)
	}
	return nil
}
