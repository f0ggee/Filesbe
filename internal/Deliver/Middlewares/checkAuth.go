package Middlewares

import (
	"Kaban/internal/Deliver/httpController"
	"Kaban/internal/DomainLevel"
	"net/http"
)

type checkAuth struct {
	TokenChecker  DomainLevel.AuthMaker
	TokenChecker2 DomainLevel.AuthMaker
	Checking      httpController.Session
}

func checkerAuth(next http.Handler, s checkAuth) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		data := s.Checking.GetSessionData(httpController.IncomingSessionData{
			Writer:  writer,
			Request: request,
		})
		if data.Error != nil {
			httpController.SetAnswer(httpController.InputAnswerData{
				W:    writer,
				Code: http.StatusUnauthorized,
				Data: httpController.AnswerUserCheck{
					UrlToRedirect: "/login",
					Error:         data.Error.Error(),
				},
			})
			return
		}
		jwtMaker, err := s.TokenChecker.Make()
		if err != nil {
			httpController.SetAnswer(httpController.InputAnswerData{
				W:    writer,
				Code: http.StatusUnauthorized,
				Data: httpController.AnswerUserCheck{
					UrlToRedirect: "/login",
					Error:         "",
				},
			})
			return
		}
		err = jwtMaker.IsTokenCorrect([]byte(data.Jwt))
		if err == nil {
			next.ServeHTTP(writer, request)
		}

		rfMaker, err := s.TokenChecker2.Make()
		if err != nil {
			httpController.SetAnswer(httpController.InputAnswerData{
				W:    writer,
				Code: http.StatusUnauthorized,
				Data: httpController.AnswerUserCheck{
					UrlToRedirect: "/login",
				},
			})
			return
		}

		err = rfMaker.IsTokenCorrect([]byte(data.Rft))
		if err != nil {
			httpController.SetAnswer(httpController.InputAnswerData{
				W:    writer,
				Code: http.StatusUnauthorized,
				Data: httpController.AnswerUserCheck{
					UrlToRedirect: "/login",
					Error:         "",
				},
			})
			return
		}
		next.ServeHTTP(writer, request)
	})
}
