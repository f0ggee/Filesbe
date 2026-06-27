package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/DeliverHandlers/Parsers"
	"Kaban/internal/Service/Application"
	"encoding/json"
	"log/slog"
	"net/http"

	"Kaban/internal/InfrastructureLayer/DeliverHandlers/SessionHandle"
)

type DeliverCases struct {
	D *Parsers.Parsing
}

func (D *DeliverCases) Register(w http.ResponseWriter, r *http.Request, s *Application.HandlerPackCollect) {

	if r.Method != http.MethodPost {
		slog.Error("Method is not POST")
		http.Error(w, "Method don't allow", http.StatusNotFound)
		return
	}

	type RegisterAnswer struct {
		StatusOfOperation string `json:"StatusOfOperation"`
		UrlRedict         string `json:"UrlRedict"`
		Error             string `json:"Error"`
	}

	userDataRegister := &Dto.UserDataRegister{}
	err := D.D.JsonParsers(userDataRegister, r)
	if err != nil {
		//TODO add handling the error
	}

	err := userDataRegister.ValidateDate()
	if err != nil {
		//TODO add handling the error
	}

	jwt, rt, err := s.RegisterService(userDataRegister, r.Context())
	//TODO add handing the error from the Register function

	returnedData := SessionHandle.SessionControl.SetNewSession(DomainLevel.IncomingSessionData{
		Writer:  w,
		Request: r,
		Jwt:     jwt,
		Rt:      rt})

	if returnedData == nil || returnedData.Error != nil {

		//TODO add handling the error
	}
	w.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(RegisterAnswer{
		StatusOfOperation: "SUCCESS",
		UrlRedict:         "/main",
	})
	if err != nil {
		slog.Error("Error is processing the json", "Error", err)
		return
	}

}
