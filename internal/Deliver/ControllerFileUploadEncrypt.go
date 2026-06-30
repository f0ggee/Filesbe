package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/SessionHandle"
	"Kaban/internal/Service/Application"
	"encoding/json"
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
		//TODO add handling the error

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
		//TODO add handling the error

		return
	}

	//TODO need to remove lines which are below
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
