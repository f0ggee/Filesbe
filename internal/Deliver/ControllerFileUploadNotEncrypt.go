package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage/AuthChecking"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoSessionHandle"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepofileUploaderNoEncryptRepo"
	"Kaban/internal/Service/Application"
	"encoding/json"

	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type FileUploaderNoEncryptNet struct {
	w http.ResponseWriter
	r *http.Request
}
type FileUploaderNoEncryptSessions struct {
	S *RepoSessionHandle.SessionConnect
}
type repoUploaderNoEncrypt struct {
	S *RepofileUploaderNoEncryptRepo.NewUploaderNoEncrypt
}

type AuthCheckingUploadEncrypt struct {
	Auth AuthChecking.NewAuthChecker
}
type UploadNotEncryptSessions struct {
	Session RepoSessionHandle.SessionConnect
}
type NewFileUploaderNoEncrypt struct {
	Net            FileUploaderNoEncryptNet
	EncryptSession FileUploaderNoEncryptNet
	Repos          repoUploaderNoEncrypt
	Auth           AuthCheckingUploadEncrypt
	Sess           UploadNotEncryptSessions
}

func (d *NewFileUploaderNoEncrypt) FileUploaderNoEncrypt(router *mux.Router, s *Application.HandlerPackCollect) {
	if r.Method != http.MethodPost {
		slog.Error("FileUploaderNoEncrypt; the method isn't allowed", slog.Group("URL details", slog.String("url", r.RequestURI)))
		http.Error(w, "err", http.StatusUnauthorized)
		return
	}
	returnedData := RepoSessionHandle.SessionControl.GetSessionData(RepoSessionHandle.IncomingSessionData{Writer: d.Net.w, Request: d.Net.r})
	if returnedData.Error != nil {
		//TODO add handling the error

		return
	}
	outData := d.Auth.Auth.CheckAuthTokens(DomainLevel.AuthCheckIncomingData{
		Jwt: returnedData.Jwt,
		Rft: returnedData.Rft,
	})
	if outData.Err != nil {
		d.Repos.S.SetBadAnswer(RepofileUploaderNoEncryptRepo.IncomingData{
			W:               d.Net.w,
			Error:           outData.Err.Error(),
			StatusOperation: DomainLevel.Break,
		})
		return
	}
	if outData.IsNewJwtCreated {
		d.Sess.Session.SetNewSession(RepoSessionHandle.IncomingSessionData{Jwt: outData.NewJwt})
	}

	filName, err := s.FileUploader(r)
	if err != nil {
		///TODO add handling the error
	}

	urlPath, err := d.Repos.S.UrlBuild(router, filName)
	if err != nil {
		d.Repos.S.SetBadAnswer(RepofileUploaderNoEncryptRepo.IncomingData{
			W:               d.Net.w,
			Error:           err.Error(),
			StatusOperation: DomainLevel.Break,
		})
		return
	}

	d.Repos.S.SetGoodAnswer(RepofileUploaderNoEncryptRepo.IncomingData{
		W:               d.Net.w,
		StatusOperation: DomainLevel.Success,
		UrlToRedirect:   urlPath,
	})
	return
}
