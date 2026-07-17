package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepofileUploaderNoEncryptRepo"
	"Kaban/internal/InfrastructureLayer/RepoSessionHandle"
	"Kaban/internal/Service/Application"
	"net/http"

	"github.com/gorilla/mux"
)

type NewFileUploaderNet struct {
	W http.ResponseWriter
	R *http.Request
}
type NewFileUploaderApp struct {
	Application.NewUpload
}
type NewFileUploaderWorkDetails struct {
	S       RepofileUploaderNoEncryptRepo.Answers
	Builder RepofileUploaderNoEncryptRepo.NewUploaderNoEncrypt
}
type NewFileUploaderSessions struct {
	Session RepoSessionHandle.Session
	Auth    AuthTokensManage.AuthCheck
}
type NewUploader struct {
	NewFileUploaderNet
	NewFileUploaderWorkDetails
	NewFileUploaderSessions
	NewFileUploaderApp
}

func GetNewFileUploader(fileUploaderNoEncryptNet NewFileUploaderNet, repoUploaderNoEncrypt NewFileUploaderWorkDetails, uploadNotEncryptSessions NewFileUploaderSessions, newFileUploaderApplication NewFileUploaderApp) *NewUploader {
	return &NewUploader{NewFileUploaderNet: fileUploaderNoEncryptNet, NewFileUploaderWorkDetails: repoUploaderNoEncrypt, NewFileUploaderSessions: uploadNotEncryptSessions, NewFileUploaderApp: newFileUploaderApplication}
}

func (d *NewUploader) FileUploaderNoEncrypt(router *mux.Router) {
	returnedData := d.Session.GetSessionData(RepoSessionHandle.IncomingSessionData{Writer: d.W, Request: d.R})
	if returnedData.Error != nil {
		d.S.SetBadAnswer(
			RepofileUploaderNoEncryptRepo.IncomingData{
				W:               d.W,
				Error:           returnedData.Error.Error(),
				StatusOperation: DomainLevel.Break,
			})
		return
	}
	outData := d.Auth.CheckUserAuth(AuthTokensManage.UserAuthCheckIncomingData{
		Jwt: returnedData.Jwt,
		Rft: returnedData.Rft,
	})
	if outData.Err != nil {
		d.S.SetBadAnswer(RepofileUploaderNoEncryptRepo.IncomingData{
			W:               d.W,
			Error:           outData.Err.Error(),
			StatusOperation: DomainLevel.Break,
		})
		return
	}
	if outData.IsNewJwtCreated {
		d.Session.SetNewSession(RepoSessionHandle.IncomingSessionData{Jwt: outData.NewJwt})
	}

	fileName, err := d.FileUploader(d.R)
	if err != nil {
		d.S.SetBadAnswer(RepofileUploaderNoEncryptRepo.IncomingData{
			W:               d.W,
			Error:           err.Error(),
			StatusOperation: DomainLevel.Break,
		})
		return
	}

	urlPath, err := d.Builder.UrlBuild(router, fileName)
	if err != nil {
		d.S.SetBadAnswer(RepofileUploaderNoEncryptRepo.IncomingData{
			W:               d.W,
			Error:           err.Error(),
			StatusOperation: DomainLevel.Break,
		})
		return
	}

	d.S.SetGoodAnswer(RepofileUploaderNoEncryptRepo.IncomingData{
		W:               d.W,
		StatusOperation: DomainLevel.Success,
		UrlToRedirect:   urlPath,
	})
	return
}
