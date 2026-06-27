package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/DeliverHandlers/SessionHandle"
	"Kaban/internal/Service/Application"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func FileUploaderEncrypt(w http.ResponseWriter, r *http.Request, router *mux.Router, s *Application.HandlerPackCollect) {

	type Answer struct {
		StatusOperation string `json:"StatusOperation"`
		Error           string `json:"Error"`
		UrlToRedict     string `json:"UrlRedict"`
	}
	if r.Method != http.MethodPost {
		slog.Error("Err in controller uploader")
		w.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
		err := json.NewEncoder(w).Encode(Answer{
			StatusOperation: DomainLevel.NotStart,
			Error:           "method don't allow",

			UrlToRedict: "nil",
		})
		if err != nil {
			return
		}

		return
	}

	returnedSessionKey := SessionHandle.SessionControl.GetSessionData(DomainLevel.IncomingSessionData{
		Writer:  w,
		Request: r,
	})
	if returnedSessionKey == nil || returnedSessionKey.Error != nil {
		//TODO add handing the error
	}
	filName, err := s.UploadEncrypt(r)
	if err != nil {
		w.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
		w.WriteHeader(400)
		if err := json.NewEncoder(w).Encode(Answer{
			StatusOperation: DomainLevel.NotStart,
			Error:           fmt.Sprint(err),
		}); err != nil {
			slog.Info("Error in encoding json ", "Error", err)
			return
		}

		return
	}

	url, err := router.Get("fileName").URL("name", filName, "bool", "true")
	if err != nil {
		slog.Error("Error can't treat", "error", err)
		return
	}

	w.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(Answer{StatusOperation: DomainLevel.Success,
		Error: "",

		UrlToRedict: url.Path}); err != nil {
		ControllerErrorLogger.ErrorContext(r.Context(), "Error in FileUploadingControlling", "Error", err)
		return
	}

}
