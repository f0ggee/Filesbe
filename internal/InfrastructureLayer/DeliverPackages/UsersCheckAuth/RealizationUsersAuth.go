package UsersCheckAuth

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"log/slog"
	"net/http"
)

type SetUsersChecker struct{}

func (s SetUsersChecker) BadAnswer(data DomainLevel.UserCheckBadIncomingData) {
	data.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	data.W.WriteHeader(http.StatusUnauthorized)
	d := DomainLevel.UserCheckAnswer{
		UrlToRedirect: data.Redirect,
		Error:         data.Err.Error(),
	}
	if err := json.NewEncoder(data.W).Encode(d); err != nil {
		slog.Error("Error decode the json", "Err", err)
		return
	}
	return

}
