package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/SessionHandle"
	"Kaban/internal/Service/Application"
	"encoding/json"
	"log/slog"
	"net/http"
)

func (d *NewRegister) Login(w http.ResponseWriter, r *http.Request, realization *Application.HandlerPackCollect) {

	type AnswerLogin struct {
		StatusOfOperation string `json:"StatusOperation"`
		UrlToRedict       string `json:"UrlRedict"`
		ErrorMessage      string `json:"ErrorMessage"`
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Dont' allow", http.StatusUnauthorized)
		slog.Error("Method Dont' allow", "Method", http.StatusUnauthorized)
		return
	}

	DataUserLogin := &Dto.UserLoginData{}
	err := d.D.JsonParsers(DataUserLogin, r)
	if err != nil {
		//TODO add handling the error

	}

	err := DataUserLogin.ValidateData()
	if err != nil {
		//TODO add handling the error
		return
	}

	JwtToken, RefreshToken, err := realization.LoginService(*DataUserLogin, r.Context())
	if err != nil {

		//TODO add handling the error
		return
	}

	ReturnedData := SessionHandle.SessionControl.SetNewSession(DomainLevel.IncomingSessionData{
		Writer:  w,
		Request: r,
		Jwt:     JwtToken,
		Rt:      RefreshToken,
	})
	if ReturnedData == nil || ReturnedData.Error != nil {
		//TODO add handling the error
	}
	w.Header().Set("Content-Type", DomainLevel.Json)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(AnswerLogin{
		StatusOfOperation: DomainLevel.Success,
		UrlToRedict:       "/main",
	}); err != nil {
		ControllerErrorLogger.ErrorContext(r.Context(), "Json in Login can't treated", "Err", err)
		return

	}

}
