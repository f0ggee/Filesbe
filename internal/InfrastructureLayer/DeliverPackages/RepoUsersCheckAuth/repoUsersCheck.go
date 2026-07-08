package RepoUsersCheckAuth

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"log/slog"
	"net/http"
)

type SetUsersChecker struct {
}
type UserAuthCheck struct {
	Jwt string
	Rft string
}
type UserCheckIncomingData struct {
	W        http.ResponseWriter
	Redirect string
	Err      error
}
type UsersCheck interface {
	BadAnswer(UserCheckIncomingData)
	SetGoodAnswer(UserCheckIncomingData)
}

func GetNewUsersChecker() *SetUsersChecker {
	return &SetUsersChecker{}
}

var UsersChecker = &SetUsersChecker{}

func (s SetUsersChecker) SetGoodAnswer(data UserCheckIncomingData) {

	data.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	data.W.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(data.W).Encode(DomainLevel.UserCheckAnswer{UrlToRedirect: data.Redirect}); err != nil {
		slog.Error("SetGoodAnswerCheckUser; error in encoding data", "ERROR", err)
		return
	}
	return
}

func (s SetUsersChecker) BadAnswer(data UserCheckIncomingData) {
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
