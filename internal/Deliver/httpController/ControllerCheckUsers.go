package httpController

import (
	"Kaban/internal/DomainLevel"
	"net/http"
)

type CheckUserAuthNetWork struct {
	writer  http.ResponseWriter
	Request *http.Request
}

func GetNewCheckUserAuthNetWork(w http.ResponseWriter, r *http.Request) *CheckUserAuthNetWork {
	return &CheckUserAuthNetWork{writer: w, Request: r}
}

type CheckUserAuth struct {
	CheckUserAuthNetWork
	Checking      Session
	TokenChecker  DomainLevel.AuthMaker
	TokenChecker2 DomainLevel.AuthMaker
}

func (s *CheckUserAuth) CheckUserAuth() {
	data := s.Checking.GetSessionData(IncomingSessionData{
		Writer:  s.writer,
		Request: s.Request,
	})
	if data.Error != nil {
		SetAnswer(InputAnswerData{
			W:    s.writer,
			Code: http.StatusUnauthorized,
			Data: AnswerUserCheck{
				UrlToRedirect: "",
				Error:         data.Error.Error(),
			},
		})
		return
	}
	jwtMaker, err := s.TokenChecker.Make()
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    s.writer,
			Code: http.StatusUnauthorized,
			Data: AnswerUserCheck{
				UrlToRedirect: "/login",
				Error:         "",
			},
		})
		return
	}
	err = jwtMaker.IsTokenCorrect([]byte(data.Jwt))
	if err == nil {
		SetAnswer(InputAnswerData{
			W:    s.writer,
			Code: http.StatusOK,
			Data: AnswerUserCheck{
				UrlToRedirect: "/main",
			},
		})
		return
	}

	rfMaker, err := s.TokenChecker2.Make()
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    s.writer,
			Code: http.StatusUnauthorized,
			Data: AnswerUserCheck{
				UrlToRedirect: "/login",
			},
		})
		return
	}

	err = rfMaker.IsTokenCorrect([]byte(data.Rft))
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    s.writer,
			Code: http.StatusUnauthorized,
			Data: AnswerUserCheck{
				UrlToRedirect: "/login",
				Error:         "",
			},
		})
		return
	}
	SetAnswer(InputAnswerData{
		W:    s.writer,
		Code: http.StatusUnauthorized,
		Data: AnswerUserCheck{
			UrlToRedirect: "/main",
			Error:         "",
		},
	})
}
