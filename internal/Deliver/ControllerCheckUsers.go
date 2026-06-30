package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/SessionHandle"
	"Kaban/internal/Service/Application"
	"encoding/json"
	"log/slog"
	"net/http"
)

func GetFrom(w http.ResponseWriter, r *http.Request, s *Application.HandlerPackCollect) {

	if r.Method != http.MethodGet {
		http.Error(w, "Cant' treat", http.StatusNotFound)
		slog.Info("Not found")
		return
	}
	type AnswerStruct struct {
		StatusRedict string `json:"status_redict"`
	}

	returnedData := SessionHandle.SessionControl.GetSessionData(DomainLevel.IncomingSessionData{Writer: w, Request: r})
	NewJwt, err := s.Auth(returnedData.Rft, returnedData.Jwt)
	if err != nil {
		w.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(AnswerStruct{StatusRedict: "/login"}); err != nil {
			slog.Error("Error decode the json", "Err", err)
			return
		}
		return
	}

	returnedData = SessionHandle.SessionControl.GetSessionData(DomainLevel.IncomingSessionData{Jwt: NewJwt})
	if returnedData == nil || returnedData.Error != nil {
		//TODO add handle the error
	}
	w.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(AnswerStruct{StatusRedict: "/main"}); err != nil {
		slog.Error("Error decode the json", "Err", err)
		return
	}
	return
}
