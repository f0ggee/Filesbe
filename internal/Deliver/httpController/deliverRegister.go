package httpController

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"Kaban/internal/InfrastructureLayer/RepoSessionHandle"
	"context"
	"net/http"
)

type NewRegisterDetails struct {
	Session RepoSessionHandle.Session
	D       RepoParsers.Decode
}

type RegisterNet struct {
	W http.ResponseWriter
	R *http.Request
}
type NewRegister struct {
	RegisterNet
	NewRegisterDetails
	RegisterService func(ctx context.Context, de *Dto.UserDataRegister) DomainLevel.RegisterApplicationOutComingData
}

func GetNewRegister(registerNet RegisterNet, newRegisterDetails NewRegisterDetails) *NewRegister {
	return &NewRegister{RegisterNet: registerNet, NewRegisterDetails: newRegisterDetails}
}

func (D NewRegister) Register() {
	userDataRegister := &Dto.UserDataRegister{}
	err := D.D.JsonDecode(userDataRegister, D.R.Body)
	defer D.R.Body.Close()

	if err != nil {
		SetAnswer(InputAnswerData{
			W:    D.W,
			code: http.StatusNotFound,
			data: AnswerRegister{
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
			code: http.StatusNotFound,
			data: AnswerRegister{
				StatusOfOperation: NotStart,
				Error:             err.Error(),
			},
		})
		return
	}

	RegisterOutput := D.RegisterService(D.R.Context(), userDataRegister)
	if RegisterOutput.Err != nil {
		SetAnswer(InputAnswerData{
			W:    D.W,
			code: http.StatusAlreadyReported,
			data: AnswerRegister{
				StatusOfOperation: Break,
				Error:             RegisterOutput.Err.Error(),
			},
		})
		return
	}
	returnedData := D.Session.SetNewSession(RepoSessionHandle.IncomingSessionData{
		Writer:  D.W,
		Request: D.R,
		Jwt:     RegisterOutput.Jwt,
		Rt:      RegisterOutput.Rft})
	if returnedData.Error != nil {
		SetAnswer(InputAnswerData{
			W:    D.W,
			code: http.StatusConflict,
			data: AnswerRegister{
				StatusOfOperation: Break,
				Error:             returnedData.Error.Error(),
			},
		})

		return
	}

	SetAnswer(InputAnswerData{
		W:    D.W,
		code: http.StatusOK,
		data: AnswerRegister{
			StatusOfOperation: Success,
			UrlToRedirect:     MainPageUrl,
		},
	})

}
