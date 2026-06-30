package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/Parsers"
	"Kaban/internal/Service/Application"
	"errors"
	"net/http"

	"Kaban/internal/InfrastructureLayer/DeliverPackages/RegisterRepo"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/SessionHandle"
)

type NewRegister struct {
	D *Parsers.Parsing
}

func (D *NewRegister) Register(w http.ResponseWriter, r *http.Request, s *Application.HandlerPackCollect) {
	userDataRegister := &Dto.UserDataRegister{}
	err := D.D.JsonParsers(userDataRegister, r.Body)
	if err != nil {
		RegisterRepo.ErrorController.ErrorAnswer(DomainLevel.RegisterErrorIncomingData{
			W:         w,
			Error:     err,
			Operation: errors.New(DomainLevel.NotStart),
		})
		return
	}

	err = userDataRegister.ValidateDate()
	if err != nil {
		RegisterRepo.ErrorController.ErrorAnswer(DomainLevel.RegisterErrorIncomingData{
			W:         w,
			Error:     err,
			Operation: errors.New(DomainLevel.NotStart),
		})
		return
	}

	jwt, rt, err := s.RegisterService(userDataRegister, r.Context())
	if err != nil {
		//TODO add handling the error
	}

	returnedData := SessionHandle.SessionControl.SetNewSession(DomainLevel.IncomingSessionData{
		Writer:  w,
		Request: r,
		Jwt:     jwt,
		Rt:      rt})
	if returnedData.Error != nil {
		RegisterRepo.ErrorController.ErrorAnswer(DomainLevel.RegisterErrorIncomingData{
			W:         w,
			Error:     returnedData.Error,
			Operation: errors.New(DomainLevel.Break),
		})
		return
	}
	RegisterRepo.ErrorController.GoodAnswer(DomainLevel.RegisterGoodIncomingData{
		W:         w,
		Operation: DomainLevel.Success,
		Redirect:  "/main",
	})

}
