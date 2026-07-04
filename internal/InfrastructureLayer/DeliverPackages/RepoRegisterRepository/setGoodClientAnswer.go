package RepoRegisterRepository

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"log/slog"
	"net/http"
)

func (r RegisterController) GoodAnswer(data RegisterGoodIncomingData) {
	d := DomainLevel.RegisterAnswer{
		StatusOfOperation: data.Operation,
		UrlToRedirect:     data.Redirect,
	}
	data.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	data.W.WriteHeader(http.StatusOK)
	err := json.NewEncoder(data.W).Encode(&d)
	if err != nil {
		slog.Error("RegisterGoodAnswer: the error happened in encoding data for a client", "ERROR", err)
		return
	}
	return
}
