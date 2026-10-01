package httpController

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/RepoSession"
	"Kaban/internal/Service/Application"
	"net/http"
)

type LoginNet struct {
	W http.ResponseWriter
	R *http.Request
}

type LoginDepends struct {
	Sess RepoSession.Session
}
type ParseLogin struct {
	Parses DomainLevel.Decoder
}

type NewLoginController struct {
	LoginNet
	LoginDepends
	Application.LoginApplication
	ParseLogin
}

func GetNewLogin(networkLogin LoginNet, loginDepends LoginDepends, parseLogin ParseLogin) *NewLoginController {
	return &NewLoginController{LoginNet: networkLogin, LoginDepends: loginDepends, ParseLogin: parseLogin}
}

func (D *NewLoginController) LoginController() {
	DataUserLogin := &Dto.UserLoginData{}
	err := D.Parses.DecodeFlow(DataUserLogin, D.R.Body)
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    D.W,
			code: http.StatusBadRequest,
			data: AnswerLogin{
				StatusOfOperation: NotStart,
				ErrorMessage:      err.Error(),
			},
		})
		return
	}

	err = DataUserLogin.ValidateData()
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    D.W,
			code: http.StatusNotFound,
			data: AnswerLogin{
				StatusOfOperation: NotStart,
				ErrorMessage:      err.Error(),
			},
		})
		return
	}

	loginDataOutput := D.Login(D.R.Context(), *DataUserLogin)
	if loginDataOutput.Err != nil {
		SetAnswer(InputAnswerData{
			W:    D.W,
			code: http.StatusUnauthorized,
			data: AnswerLogin{
				StatusOfOperation: Break,
				ErrorMessage:      loginDataOutput.Err.Error(),
			},
		})
		return
	}
	ReturnedData := D.Sess.SetNewSession(RepoSession.IncomingSessionData{
		Writer:  D.W,
		Request: D.R,
		Jwt:     loginDataOutput.Jwt,
		Rt:      loginDataOutput.Rft,
	})
	if ReturnedData.Error != nil {
		SetAnswer(InputAnswerData{
			W:    D.W,
			code: http.StatusBadRequest,
			data: AnswerLogin{
				StatusOfOperation: Break,
				ErrorMessage:      ReturnedData.Error.Error(),
			},
		})
		return
	}
	SetAnswer(InputAnswerData{
		W:    D.W,
		code: http.StatusOK,
		data: AnswerLogin{
			StatusOfOperation: Success,
			UrlToRedirect:     MainPageUrl,
		},
	})

	return
}
