package httpController

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/Service/Application"
	"net/http"
)

type NewLoginController struct {
	Sess   Session
	Parses DomainLevel.Decoder
	W      http.ResponseWriter
	R      *http.Request
	App    Application.LoginApplication
}

func (D *NewLoginController) LoginController() {
	DataUserLogin := &Dto.UserLoginData{}
	err := D.Parses.DecodeFlow(DataUserLogin, D.R.Body)
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    D.W,
			Code: http.StatusBadRequest,
			Data: AnswerLogin{
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
			Code: http.StatusNotFound,
			Data: AnswerLogin{
				StatusOfOperation: NotStart,
				ErrorMessage:      err.Error(),
			},
		})
		return
	}

	loginDataOutput := D.App.Login(D.R.Context(), *DataUserLogin)
	if loginDataOutput.Err != nil {
		SetAnswer(InputAnswerData{
			W:    D.W,
			Code: http.StatusUnauthorized,
			Data: AnswerLogin{
				StatusOfOperation: Break,
				ErrorMessage:      loginDataOutput.Err.Error(),
			},
		})
		return
	}
	ReturnedData := D.Sess.SetNewSession(IncomingSessionData{
		Writer:  D.W,
		Request: D.R,
		Jwt:     loginDataOutput.Jwt,
		Rt:      loginDataOutput.Rft,
	})
	if ReturnedData.Error != nil {
		SetAnswer(InputAnswerData{
			W:    D.W,
			Code: http.StatusBadRequest,
			Data: AnswerLogin{
				StatusOfOperation: Break,
				ErrorMessage:      ReturnedData.Error.Error(),
			},
		})
		return
	}
	SetAnswer(InputAnswerData{
		W:    D.W,
		Code: http.StatusOK,
		Data: AnswerLogin{
			StatusOfOperation: Success,
			UrlToRedirect:     MainPageUrl,
		},
	})

	return
}
