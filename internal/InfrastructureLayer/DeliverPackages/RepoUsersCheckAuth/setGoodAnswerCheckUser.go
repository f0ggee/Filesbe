package RepoUsersCheckAuth

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"log/slog"
	"net/http"
)

func (s SetUsersChecker) SetGoodAnswer(data UserCheckIncomingData) {

	data.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	data.W.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(data.W).Encode(DomainLevel.UserCheckAnswer{UrlToRedirect: data.Redirect}); err != nil {
		slog.Error("SetGoodAnswerCheckUser; error in encoding data", "ERROR", err)
		return
	}
	return
}
