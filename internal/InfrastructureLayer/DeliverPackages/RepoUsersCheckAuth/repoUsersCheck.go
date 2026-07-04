package RepoUsersCheckAuth

import "net/http"

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
