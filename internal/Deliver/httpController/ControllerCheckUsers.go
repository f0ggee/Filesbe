package httpController

import (
	"net/http"
)

type CheckUserAuthNetWork struct {
	writer  http.ResponseWriter
	Request *http.Request
}

func GetNewCheckUserAuthNetWork(w http.ResponseWriter, r *http.Request) *CheckUserAuthNetWork {
	return &CheckUserAuthNetWork{writer: w, Request: r}
}

type CheckUserAuthController struct {
	CheckUserAuthNetWork
}

func (s *CheckUserAuthController) CheckUserAuth() {
	SetAnswer(InputAnswerData{
		W:    s.writer,
		Code: http.StatusUnauthorized,
		Data: AnswerUserCheck{
			UrlToRedirect: "/main",
			Error:         "",
		},
	})
}
