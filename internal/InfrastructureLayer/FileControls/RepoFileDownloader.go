package FileControls

import (
	"Kaban/internal/DomainLevel"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type FileDetails struct {
	FileFormat   string
	TrueFileName string
	FileLength   int64
}
type TransferIncomingData struct {
	W http.ResponseWriter
	FileDetails
	FileBody io.ReadCloser
}

type Transferring interface {
	TransferToClient(TransferIncomingData) error
	TransferEncryptToClient(TransferIncomingData) error
}

type Transfer struct{}

func GetNewTransfer() *Transfer {
	return &Transfer{}
}

func (n Transfer) TransferEncryptToClient(details TransferIncomingData) error {
	details.W.Header().Set("Content-Type", details.FileFormat)
	details.W.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename= %v", details.TrueFileName))
	details.W.Header().Set("Content-Length", strconv.FormatUint(uint64(details.FileLength), 10))
	if _, err := io.Copy(details.W, details.FileBody); err != nil {
		slog.Error("TransferEncryptToClient; error to download a file", "ERROR", err)
		return errors.New(DomainLevel.ErrorDownloadFile)
	}
	return nil
}
func (n Transfer) TransferToClient(details TransferIncomingData) error {
	details.W.Header().Set("Content-Type", details.FileFormat)
	details.W.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename= %v", details.TrueFileName))
	details.W.Header().Set("Content-Length", strconv.FormatUint(uint64(details.FileLength), 10))
	if _, err := io.Copy(details.W, details.FileBody); err != nil {
		slog.Error("TransferToClient; error to download a file", "ERROR", err)
		return errors.New(DomainLevel.ErrorDownloadFile)
	}
	return nil
}
