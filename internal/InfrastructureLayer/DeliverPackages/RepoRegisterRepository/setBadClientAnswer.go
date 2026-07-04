package RepoRegisterRepository

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"log/slog"
	"net/http"
)

func (r RegisterController) ErrorAnswer(data RegisterErrorIncomingData) {
	userErrorResponse := DomainLevel.RegisterAnswer{
		StatusOfOperation: data.Operation.Error(),
		UrlToRedirect:     "",
		Error:             data.Error.Error(),
	}
	data.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	data.W.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(data.W).Encode(userErrorResponse); err != nil {
		slog.Error("ErrorAnswerRegister; error during encoding", "ERROR", err)
		return
	}
	return
}
