package httpController

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type InputAnswerData struct {
	W    http.ResponseWriter
	Code int
	Data any
}

func SetAnswer(data InputAnswerData) {
	data.W.Header().Set("Content-Type", "application/json")
	data.W.WriteHeader(data.Code)
	if err := json.NewEncoder(data.W).Encode(&data.Data); err != nil {
		slog.Error("SetAnswer; error to encode a response", "ERROR", err)
		return
	}
	return

}
