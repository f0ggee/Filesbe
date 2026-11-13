package httpController

import (
	"Kaban/internal/Deliver/httpController/httpData"
	"Kaban/internal/DomainLevel"
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
	var dataRegister httpData.RegisterData

	err := D.Decoder.DecodeFlow(&dataRegister, D.R.Body)
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

	err = dataRegister.ValidateDate()
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

	RegisterOutput := D.App.RegisterApplication(Application.InputRegister{
		Ctx:      D.R.Context(),
		Name:     dataRegister.Name,
		Email:    dataRegister.Email,
		Password: dataRegister.Password,
	})
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
