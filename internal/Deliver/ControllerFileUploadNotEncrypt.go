package Deliver

import "C"
import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoSessionHandle"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepofileUploaderNoEncryptRepo"
	"Kaban/internal/Service/Application"
	"net/http"

	"github.com/gorilla/mux"
)

type FileUploaderNoEncryptNet struct {
	w http.ResponseWriter
	r *http.Request
}
type NewFileUploaderApplication struct {
	Application.NewFileUploader
}
type FileUploaderNoEncryptSessions struct {
	S *RepoSessionHandle.SessionConnect
}
type repoUploaderNoEncrypt struct {
	S *RepofileUploaderNoEncryptRepo.NewUploaderNoEncrypt
}

type AuthCheckingUploadEncrypt struct {
	Auth AuthTokensManage.NewAuthChecker
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
	App            NewFileUploaderApplication
}

func GetNewFileUploaderNoEncrypt(net FileUploaderNoEncryptNet, encryptSession FileUploaderNoEncryptNet, repos repoUploaderNoEncrypt, auth AuthCheckingUploadEncrypt, sess UploadNotEncryptSessions, app NewFileUploaderApplication) *NewFileUploaderNoEncrypt {
	return &NewFileUploaderNoEncrypt{Net: net, EncryptSession: encryptSession, Repos: repos, Auth: auth, Sess: sess, App: app}
}

func (d *NewFileUploaderNoEncrypt) FileUploaderNoEncrypt(router *mux.Router) {
	//TODO here is the http Post
	returnedData := RepoSessionHandle.SessionControl.GetSessionData(RepoSessionHandle.IncomingSessionData{Writer: d.Net.w, Request: d.Net.r})
	if returnedData.Error != nil {
		//TODO add handling the error

		return
	}
	outData := d.Auth.Auth.CheckUserAuth(AuthTokensManage.UserAuthCheckIncomingData{
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

	fileName, err := d.App.FileUploader(d.Net.r)
	if err != nil {
		d.Repos.S.SetBadAnswer(RepofileUploaderNoEncryptRepo.IncomingData{
			W:               d.Net.w,
			Error:           err.Error(),
			StatusOperation: DomainLevel.Break,
		})
		return
	}

	urlPath, err := d.Repos.S.UrlBuild(router, fileName)
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
