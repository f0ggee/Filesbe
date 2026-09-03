package httpController

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"Kaban/internal/InfrastructureLayer/RepoSessionHandle"
	"context"

	"net/http"
)

type LoginNet struct {
	W http.ResponseWriter
	R *http.Request
}

type LoginDepends struct {
	Sess RepoSessionHandle.Session
}
type ParseLogin struct {
	Parses RepoParsers.Decode
}

type NewLoginController struct {
	LoginNet
	LoginDepends
	LoginApp func(ctx context.Context, s Dto.UserLoginData) DomainLevel.LoginApplicationOutComingData
	ParseLogin
}

func GetNewLogin(networkLogin LoginNet, loginDepends LoginDepends, parseLogin ParseLogin) *NewLoginController {
	return &NewLoginController{LoginNet: networkLogin, LoginDepends: loginDepends, ParseLogin: parseLogin}
}

func (D *NewLoginController) Login() {
	DataUserLogin := &Dto.UserLoginData{}
	err := D.Parses.JsonDecode(DataUserLogin, D.R.Body)
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

	loginDataOutput := D.LoginApp(D.R.Context(), *DataUserLogin)
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
	ReturnedData := D.Sess.SetNewSession(RepoSessionHandle.IncomingSessionData{
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
