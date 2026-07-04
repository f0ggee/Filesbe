package RepoLoginRealizations

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"log/slog"
	"net/http"
)

type LoginAnswers struct{}

func (l LoginAnswers) SetGoodAnswer(data AnswerDetails) {
	data.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	data.W.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(data.W).Encode(DomainLevel.AnswerLogin{
		StatusOfOperation: DomainLevel.Success,
		UrlToRedirect:     "/main",
	}); err != nil {
		slog.Error("SetLoginGoodAnswer; error during encoding", "ERROR", err)
		return
	}
}

func (l LoginAnswers) SetBadAnswer(details AnswerDetails) {
	details.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	details.W.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(details.W).Encode(DomainLevel.AnswerLogin{
		StatusOfOperation: details.Operation,
		ErrorMessage:      details.Error,
	}); err != nil {
		slog.Error("SetLoginBadAnswer; error during encoding", "ERROR", err)
		return
	}
}

type AnswerDetails struct {
	W               http.ResponseWriter
	Operation       string
	Error           string
	UrlToRedistrict string
}

type TypeAnswers interface {
	SetBadAnswer(AnswerDetails)
	SetGoodAnswer(AnswerDetails)
}
