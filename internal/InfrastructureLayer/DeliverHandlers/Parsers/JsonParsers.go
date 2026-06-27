package Parsers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

func (p Parsing) JsonParsers(a any, request *http.Request) error {
	if err := json.NewDecoder(request.Body).Decode(&a); err != nil {
		slog.Error("GetIncomingData: the error to parse data")
		return nil

	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Error is closing the body in the controller register", "ERROR", err)
			return
		}
	}(request.Body)

	return nil
}
