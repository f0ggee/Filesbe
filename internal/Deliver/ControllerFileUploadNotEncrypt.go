package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/SessionHandle"
	"Kaban/internal/Service/Application"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func FileUploaderNoEncrypt(w http.ResponseWriter, r *http.Request, router *mux.Router, s *Application.HandlerPackCollect) {
	if r.Method != http.MethodPost {
		http.Error(w, "err", http.StatusUnauthorized)
		slog.Error("Err in Cottroler Uploader")
		return
	}
	type Answer struct {
		StatusOperation string `json:"StatusOperation"`
		UrlToRedict     string `json:"UrlRedict"`
		Error           string `json:"Error"`
	}

	returnedData := SessionHandle.SessionControl.GetSessionData(DomainLevel.IncomingSessionData{Writer: w, Request: r})
	if returnedData == nil || returnedData.Error != nil {
		//TODO add handling the error

		return
	}
	Jwts, err := s.Auth(returnedData.Rft, returnedData.Jwt)
	if err != nil {
		//TODO add handling the error
		return
	}

	filName, err := s.FileUploader(r)
	if err != nil {
		///TODO add handling the error
	}

	url, err := router.Get("fileName").URL("name", filName, "bool", "false")
	if err != nil {
		slog.Error("Error can't treate", "Error", err)

		w.Header().Set("Content-Type", DomainLevel.Json)
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(Answer{
			StatusOperation: DomainLevel.Break,
		}); err != nil {
			slog.Error("Err in json encode", "error", err)
			return
		}
		return
	}

	w.Header().Set("Content-Type", DomainLevel.Json)
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(Answer{
		StatusOperation: DomainLevel.Success,
		UrlToRedict:     url.Path,
	}); err != nil {
		slog.Error("Err in json encode", "Error", err)
		return
	}

}
