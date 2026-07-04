package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoParsers"
	"errors"
	"net/http"

	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoRegisterRepository"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoSessionHandle"
)

type NewRegister struct {
	D *RepoParsers.Parsing
	W http.ResponseWriter
	R *http.Request
	S *RepoSessionHandle.Session
}

func (D NewRegister) Register() {
	userDataRegister := &Dto.UserDataRegister{}
	err := D.D.JsonDecode(userDataRegister, D.R.Body)
	if err != nil {
		RepoRegisterRepository.ErrorController.ErrorAnswer(RepoRegisterRepository.RegisterErrorIncomingData{
			W:         D.W,
			Error:     err,
			Operation: errors.New(DomainLevel.NotStart),
		})
		return
	}

	err = userDataRegister.ValidateDate()
	if err != nil {
		RepoRegisterRepository.ErrorController.ErrorAnswer(RepoRegisterRepository.RegisterErrorIncomingData{
			W:         D.W,
			Error:     err,
			Operation: errors.New(DomainLevel.NotStart),
		})
		return
	}

	jwt, rt, err := D.S.RegisterService(userDataRegister, D.R.Context())
	if err != nil {
		//TODO add handling the error
	}

	returnedData := RepoSessionHandle.SessionControl.SetNewSession(RepoSessionHandle.IncomingSessionData{
		Writer:  D.W,
		Request: D.R,
		Jwt:     jwt,
		Rt:      rt})
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
