package RepoParsers

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
)

func (p Parsing) JsonDecode(a any, request io.ReadCloser) error {
	defer func() {
		err := request.Close()
		if err != nil {
			slog.Error("JsonDecode: the error happened", "ERROR", err)
			return
		}
	}()
	if err := json.NewDecoder(request).Decode(&a); err != nil {
		slog.Error("GetIncomingData: the error to parse data", "ERROR", err)
		return errors.New(DomainLevel.ErrorParseInfo)
	}
	return nil
}
