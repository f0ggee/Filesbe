package httpController

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type InputAnswerData struct {
	W    http.ResponseWriter
	code int
	data any
}

func SetAnswer(data InputAnswerData) {
	data.W.Header().Set(ContentType, Json)
	data.W.WriteHeader(data.code)
	if err := json.NewEncoder(data.W).Encode(&data.data); err != nil {
		slog.Error("SetAnswer; error to encode a response", "ERROR", err)
		return
	}
	return

}
