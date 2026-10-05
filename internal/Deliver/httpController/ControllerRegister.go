package httpController

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/Service/Application"
	"net/http"
)

type RegisterNet struct {
	W http.ResponseWriter
	R *http.Request
}
type RegisterController struct {
	RegisterNet
	Session Session
	Decoder DomainLevel.Decoder
	App     Application.RegisterApplication
}

func (D RegisterController) Register() {
	var userDataRegister Dto.UserDataRegister
	err := D.Decoder.DecodeFlow(&userDataRegister, D.R.Body)
	defer D.R.Body.Close()

	if err != nil {
		SetAnswer(InputAnswerData{
			W:    D.W,
			Code: http.StatusNotFound,
			Data: AnswerRegister{
				StatusOfOperation: DomainName,
				Error:             err.Error(),
			},
		})
		return
	}

	err = userDataRegister.ValidateDate()
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    D.W,
			Code: http.StatusNotFound,
			Data: AnswerRegister{
				StatusOfOperation: NotStart,
				Error:             err.Error(),
			},
		})
		return
	}

	RegisterOutput := D.App.RegisterApp(D.R.Context(), &userDataRegister)
	if RegisterOutput.Err != nil {
		SetAnswer(InputAnswerData{
			W:    D.W,
			Code: http.StatusAlreadyReported,
			Data: AnswerRegister{
				StatusOfOperation: Break,
				Error:             RegisterOutput.Err.Error(),
			},
		})
		return
	}
	returnedData := D.Session.SetNewSession(IncomingSessionData{
		Writer:  D.W,
		Request: D.R,
		Jwt:     RegisterOutput.Jwt,
		Rt:      RegisterOutput.Rft})
	if returnedData.Error != nil {
		SetAnswer(InputAnswerData{
			W:    D.W,
			Code: http.StatusConflict,
			Data: AnswerRegister{
				StatusOfOperation: Break,
				Error:             returnedData.Error.Error(),
			},
		})

		return
	}

	SetAnswer(InputAnswerData{
		W:    D.W,
		Code: http.StatusOK,
		Data: AnswerRegister{
			StatusOfOperation: Success,
			UrlToRedirect:     MainPageUrl,
		},
	})

}
