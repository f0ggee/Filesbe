package httpController

import (
	"Kaban/internal/Deliver/httpController/httpData"
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Service/Application"
	"net/http"
)

type LoginController struct {
	Sess   Session
	Parses DomainLevel.Decoder
	App    Application.LoginApplication
}

func (D *LoginController) LoginController(W http.ResponseWriter, R *http.Request) {
	var DataUserLogin httpData.LoginData
	err := D.Parses.DecodeFlow(&DataUserLogin, R.Body)
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    W,
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
			W:    W,
			Code: http.StatusNotFound,
			Data: AnswerLogin{
				StatusOfOperation: NotStart,
				ErrorMessage:      err.Error(),
			},
		})
		return
	}

	loginDataOutput := D.App.LoginApplication(Application.InputLogin{
		Ctx:      R.Context(),
		Email:    DataUserLogin.Email,
		Password: DataUserLogin.Password,
	})
	if loginDataOutput.Err != nil {
		SetAnswer(InputAnswerData{
			W:    W,
			Code: http.StatusUnauthorized,
			Data: AnswerLogin{
				StatusOfOperation: Break,
				ErrorMessage:      loginDataOutput.Err.Error(),
			},
		})
		return
	}
	ReturnedData := D.Sess.SetNewSession(IncomingSessionData{
		Writer:  W,
		Request: R,
		Jwt:     loginDataOutput.Jwt,
		Rt:      loginDataOutput.Rft,
	})
	if ReturnedData.Error != nil {
		SetAnswer(InputAnswerData{
			W:    W,
			Code: http.StatusBadRequest,
			Data: AnswerLogin{
				StatusOfOperation: Break,
				ErrorMessage:      ReturnedData.Error.Error(),
			},
		})
		return
	}
	SetAnswer(InputAnswerData{
		W:    W,
		Code: http.StatusOK,
		Data: AnswerLogin{
			StatusOfOperation: Success,
			UrlToRedirect:     MainPageUrl,
		},
	})

	return
}
