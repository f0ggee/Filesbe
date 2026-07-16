package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"Kaban/internal/InfrastructureLayer/RepoSessionHandle"
	"Kaban/internal/Service/Application"
	"errors"
	"net/http"

	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoRegisterRepository"
)

type NewRegisterDetails struct {
	Answ    RepoRegisterRepository.RegisterAnswers
	Session RepoSessionHandle.Session
	D       RepoParsers.Decode
}
type NewRegisterApp struct {
	App Application.NewRegisterApplication
}
type RegisterNet struct {
	W http.ResponseWriter
	R *http.Request
}
type NewRegister struct {
	RegisterNet
	NewRegisterDetails
	NewRegisterApp
}

func GetNewRegister(registerNet RegisterNet, newRegisterDetails NewRegisterDetails, newRegisterApp NewRegisterApp) *NewRegister {
	return &NewRegister{RegisterNet: registerNet, NewRegisterDetails: newRegisterDetails, NewRegisterApp: newRegisterApp}
}

func (D NewRegister) Register() {
	userDataRegister := &Dto.UserDataRegister{}
	err := D.D.JsonDecode(userDataRegister, D.R.Body)
	if err != nil {
		D.Answ.ErrorAnswer(RepoRegisterRepository.RegisterErrorIncomingData{
			W:         D.W,
			Error:     err,
			Operation: errors.New(DomainLevel.NotStart),
		})
		return
	}

	err = userDataRegister.ValidateDate()
	if err != nil {
		D.Answ.ErrorAnswer(RepoRegisterRepository.RegisterErrorIncomingData{
			W:         D.W,
			Error:     err,
			Operation: errors.New(DomainLevel.NotStart),
		})
		return
	}

	RegisterOutput := D.App.RegisterService(userDataRegister, D.R.Context())
	if RegisterOutput.Err != nil {
		D.NewRegisterDetails.Answ.ErrorAnswer(RepoRegisterRepository.RegisterErrorIncomingData{
			W:         D.W,
			Error:     RegisterOutput.Err,
			Operation: errors.New(DomainLevel.Break),
		})
		return
	}
	returnedData := D.Session.SetNewSession(RepoSessionHandle.IncomingSessionData{
		Writer:  D.W,
		Request: D.R,
		Jwt:     RegisterOutput.Jwt,
		Rt:      RegisterOutput.Rft})
	if returnedData.Error != nil {
		RepoRegisterRepository.ErrorController.ErrorAnswer(RepoRegisterRepository.RegisterErrorIncomingData{
			W:         D.W,
			Error:     returnedData.Error,
			Operation: errors.New(DomainLevel.Break),
		})
		return
	}
	RepoRegisterRepository.ErrorController.GoodAnswer(RepoRegisterRepository.RegisterGoodIncomingData{
		W:         D.W,
		Operation: DomainLevel.Success,
		Redirect:  "/main",
	})

}
