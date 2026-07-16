package Deliver

import (
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoUsersCheckAuth"
	"Kaban/internal/InfrastructureLayer/RepoSessionHandle"
	"log/slog"
	"net/http"
)

type NewCheckUserAuthNetWork struct {
	W http.ResponseWriter
	R *http.Request
}

func GetNewCheckUserAuthNetWork(w http.ResponseWriter, r *http.Request) *NewCheckUserAuthNetWork {
	return &NewCheckUserAuthNetWork{W: w, R: r}
}

type NewCheckUserAuthWorkDetails struct {
	Answers *RepoUsersCheckAuth.SetUsersChecker
}

func GetNewCheckUserAuthWorkDetails(answers *RepoUsersCheckAuth.SetUsersChecker) *NewCheckUserAuthWorkDetails {
	return &NewCheckUserAuthWorkDetails{Answers: answers}
}

type NewCheckUserAuthSessions struct {
	Session RepoSessionHandle.Session
	Auth    AuthTokensManage.AuthCheck
}

func GetNewCheckUserAuthSessions(auth AuthTokensManage.AuthCheck, session RepoSessionHandle.Session) *NewCheckUserAuthSessions {
	return &NewCheckUserAuthSessions{Auth: auth, Session: session}
}

type NewCheckUserAuth struct {
	NewCheckUserAuthNetWork
	NewCheckUserAuthWorkDetails
	NewCheckUserAuthSessions
}

func GetNewCheckUserAuth(netWork NewCheckUserAuthNetWork, checkUserAuthDetails NewCheckUserAuthWorkDetails, sessions NewCheckUserAuthSessions) *NewCheckUserAuth {
	return &NewCheckUserAuth{NewCheckUserAuthNetWork: netWork, NewCheckUserAuthWorkDetails: checkUserAuthDetails, NewCheckUserAuthSessions: sessions}
}
func (s *NewCheckUserAuth) CheckUserAuth() {
	//The post method is here
	if s.R.Method != http.MethodGet {
		slog.Error("CheckUserAuth; Method isn't allowed", slog.Group("Details", slog.String("Method", s.R.Method), slog.String("The url", s.R.RequestURI)))
		return
	}

	returnedData := s.Session.GetSessionData(RepoSessionHandle.IncomingSessionData{Writer: s.W, Request: s.R})
	if returnedData.Error != nil {
		s.Answers.BadAnswer(RepoUsersCheckAuth.UserCheckIncomingData{
			W:   s.W,
			Err: returnedData.Error,
		})
		return
	}
	OutData := s.Auth.CheckUserAuth(AuthTokensManage.UserAuthCheckIncomingData{
		Jwt: returnedData.Jwt,
		Rft: returnedData.Rft,
	})
	if OutData.Err != nil {
		RepoUsersCheckAuth.UsersChecker.BadAnswer(RepoUsersCheckAuth.UserCheckIncomingData{
			W:        s.W,
			Redirect: "/login",
			Err:      returnedData.Error,
		})
		return
	}
	if OutData.IsNewJwtCreated {
		s.Session.SetNewSession(RepoSessionHandle.IncomingSessionData{
			Writer:  s.W,
			Request: s.R,
			Jwt:     OutData.NewJwt,
		})
	}
	s.Answers.SetGoodAnswer(RepoUsersCheckAuth.UserCheckIncomingData{W: s.W, Redirect: "/main"})
	return
}
