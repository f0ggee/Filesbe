package TransferringHttpRepo

import (
	"Kaban/internal/DomainLevel"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type HttpUploader struct {
	sizeFile int64
	name     string
	source   io.Writer
}

func NewHttpUploader() HttpUploader {
	return HttpUploader{}
}

func (h *HttpUploader) SetName(s string) DomainLevel.MakerUploader {
	h.name = s
	return h
}

func (h *HttpUploader) SetSize(i int64) DomainLevel.MakerUploader {

	h.sizeFile = i
	return h
}

func (h *HttpUploader) SetAdditionalWriter(writer io.Writer) DomainLevel.MakerUploader {
	h.source = writer
	return h
}

func (h *HttpUploader) Uploader(reader io.Reader) error {
	wri := h.source.(http.ResponseWriter)
	wri.Header().Set("Content-Type", DomainLevel.GetNewFileSettings(0, h.name).FindFormatOfFile())
	wri.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename= %v", h.name))
	wri.Header().Set("Content-Length", strconv.FormatUint(uint64(h.sizeFile), 10))
	if _, err := io.Copy(h.source, reader); err != nil {
		slog.Error("HttpUploader: error to upload a file", "ERROR", err)
		return ErrorUploadFile
	}
	return nil
}

func (h *HttpUploader) CloseSource() error {

	return nil
}

func (h *HttpUploader) GetFileSize() int64 {
	return h.sizeFile
}

func (h *HttpUploader) GetFileName() string {
	return h.name
}

func (h *HttpUploader) Make(ctx context.Context) (DomainLevel.Upload, error) {
	return h, nil
}
