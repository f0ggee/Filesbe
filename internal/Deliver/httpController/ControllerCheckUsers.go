package httpController

import (
	"net/http"
)

type CheckUserAuthController struct {
}

func (s *CheckUserAuthController) CheckUserAuth(W http.ResponseWriter, R *http.Request) {
	SetAnswer(InputAnswerData{
		W:    W,
		Code: http.StatusUnauthorized,
		Data: AnswerUserCheck{
			UrlToRedirect: "/main",
			Error:         "",
		},
	})
}
