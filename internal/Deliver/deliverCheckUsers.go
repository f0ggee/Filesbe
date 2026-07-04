package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage/AuthChecking"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage/Creating"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage/ValidatingTokens"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoSessionHandle"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoUsersCheckAuth"
	"log/slog"
	"net/http"
)

type NetWork struct {
	W http.ResponseWriter
	R *http.Request
}
type TokensChecker struct {
	TokenCreate *Creating.CreatingTokens
	TokenCheck  *ValidatingTokens.Checking
}

type CheckAuth struct {
	Answers *RepoUsersCheckAuth.SetUsersChecker
}
type Sessions struct {
	Session *RepoSessionHandle.SessionConnect
}
type AuthChecker struct {
	Auth AuthChecking.NewAuthChecker
}
type NewCheckUserAuth struct {
	Net     NetWork
	Tokens  TokensChecker
	Answ    CheckAuth
	Session Sessions
	NewAuth AuthChecker
}

func (s *NewCheckUserAuth) CheckUserAuth() {
	if s.Net.R.Method != http.MethodGet {
		slog.Error("CheckUserAuth; Method isn't allowed", slog.Group("Details", slog.String("Method", s.Net.R.Method), slog.String("The url", s.Net.R.RequestURI)))
		return
	}

	returnedData := RepoSessionHandle.SessionControl.GetSessionData(RepoSessionHandle.IncomingSessionData{Writer: s.Net.W, Request: s.Net.R})
	if returnedData.Error != nil {
		s.Answ.Answers.BadAnswer(RepoUsersCheckAuth.UserCheckIncomingData{
			W:   s.Net.W,
			Err: returnedData.Error,
		})
		return
	}
	OutData := s.NewAuth.Auth.CheckAuthTokens(DomainLevel.AuthCheckIncomingData{
		Jwt: returnedData.Jwt,
		Rft: returnedData.Rft,
	})
	if OutData.Err != nil {
		RepoUsersCheckAuth.UsersChecker.BadAnswer(RepoUsersCheckAuth.UserCheckIncomingData{
			W:        s.Net.W,
			Redirect: "/login",
			Err:      returnedData.Error,
		})
		return
	}
	if OutData.IsNewJwtCreated {
		s.Session.Session.SetNewSession(RepoSessionHandle.IncomingSessionData{
			Writer:  s.Net.W,
			Request: s.Net.R,
			Jwt:     OutData.NewJwt,
		})
	}
	s.Answ.Answers.SetGoodAnswer(RepoUsersCheckAuth.UserCheckIncomingData{W: s.Net.W, Redirect: "/main"})
	return
}
