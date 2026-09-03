package httpController

import (
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"Kaban/internal/InfrastructureLayer/RepoSessionHandle"
	"net/http"
)

type CheckUserAuthNetWork struct {
	W http.ResponseWriter
	R *http.Request
}

func GetNewCheckUserAuthNetWork(w http.ResponseWriter, r *http.Request) *CheckUserAuthNetWork {
	return &CheckUserAuthNetWork{W: w, R: r}
}

type CheckUserAuthSessions struct {
	Session RepoSessionHandle.Session
	Auth    AuthTokensManage.AuthCheck
}

func GetNewCheckUserAuthSessions(auth AuthTokensManage.AuthCheck, session RepoSessionHandle.Session) *CheckUserAuthSessions {
	return &CheckUserAuthSessions{Auth: auth, Session: session}
}

type CheckUserAuth struct {
	CheckUserAuthNetWork
	CheckUserAuthSessions
}

func GetNewCheckUserAuth(newCheckUserAuthNetWork CheckUserAuthNetWork, newCheckUserAuthSessions CheckUserAuthSessions) *CheckUserAuth {
	return &CheckUserAuth{CheckUserAuthNetWork: newCheckUserAuthNetWork, CheckUserAuthSessions: newCheckUserAuthSessions}
}

func (s *CheckUserAuth) CheckUserAuth() {
	returnedData := s.Session.GetSessionData(RepoSessionHandle.IncomingSessionData{Writer: s.W, Request: s.R})
	if returnedData.Error != nil {
		SetAnswer(InputAnswerData{
			W:    s.W,
			code: http.StatusUnauthorized,
			data: AnswerUserCheck{
				Error: returnedData.Error.Error(),
			},
		})
		return
	}
	OutData := s.Auth.CheckUserAuth(AuthTokensManage.UserAuthCheckIncomingData{
		Jwt: returnedData.Jwt,
		Rft: returnedData.Rft,
	})
	if OutData.Err != nil {
		SetAnswer(InputAnswerData{
			W:    s.W,
			code: http.StatusUnauthorized,
			data: AnswerUserCheck{
				UrlToRedirect: LoginPage,
				Error:         returnedData.Error.Error(),
			},
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

	SetAnswer(InputAnswerData{
		W:    s.W,
		code: http.StatusOK,
		data: AnswerUserCheck{
			UrlToRedirect: MainPageUrl,
		},
	})
	return
}
