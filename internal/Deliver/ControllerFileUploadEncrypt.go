package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoSessionHandle"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepofileUploaderEncryptRepo"
	"Kaban/internal/Service/Application"
	"net/http"

	"github.com/gorilla/mux"
)

type FileUploaderEncryptNetwork struct {
	W http.ResponseWriter
	R *http.Request
}
type NewFileUploaderEncryptApplication struct {
	Application.NewUploadEncrypt
}

type FileUploaderEncryptRouting struct {
	Rout *mux.Router
}
type Session struct {
	ReadSession RepoSessionHandle.NewSessionConnect
	AuthCheck   AuthTokensManage.AuthCheck
}
type AnswerUploadEncrypt struct {
	S RepofileUploaderEncryptRepo.SetNewUploadingRepo
}
type NewFileUploaderEncrypt struct {
	Net         FileUploaderNoEncryptNet
	Router      FileUploaderEncryptRouting
	Sess        Session
	Answe       AnswerUploadEncrypt
	Application NewFileUploaderEncryptApplication
}

func (S *NewFileUploaderEncrypt) FileUploaderEncrypt() {

	returnedSession := S.Sess.ReadSession.GetSessionData(RepoSessionHandle.IncomingSessionData{
		Writer:  S.Net.w,
		Request: S.Net.r,
	})
	if returnedSession.Error != nil {
		S.Answe.S.SetBadAnswers(RepofileUploaderEncryptRepo.IncomingDataAnswer{
			W:               S.Net.w,
			Error:           returnedSession.Error.Error(),
			StatusOperation: DomainLevel.Break,
		})
		return
	}
	Data := S.Sess.AuthCheck.CheckUserAuth(AuthTokensManage.UserAuthCheckIncomingData{
		Jwt: returnedSession.Jwt,
		Rft: returnedSession.Rft,
	})
	if Data.Err != nil {
		S.Answe.S.SetGoodAnswers(RepofileUploaderEncryptRepo.IncomingDataAnswer{
			W:               S.Net.w,
			Error:           Data.Err.Error(),
			StatusOperation: DomainLevel.Break,
		})
		return
	}

	fileName, err := S.Application.NewUploadEncrypt.UploadEncrypt(S.Net.r)
	if err != nil {
		//TODO add handling the error

		return
	}

	urlPath, err := S.Answe.S.UrlBuilder(S.Router.Rout, fileName)
	if err != nil {
		S.Answe.S.SetBadAnswers(RepofileUploaderEncryptRepo.IncomingDataAnswer{
			W:               S.Net.w,
			Error:           err.Error(),
			StatusOperation: DomainLevel.NotStart,
		})
		return
	}

	S.Answe.S.SetGoodAnswers(RepofileUploaderEncryptRepo.IncomingDataAnswer{
		W:               S.Net.w,
		UrlToRedirect:   urlPath,
		StatusOperation: DomainLevel.Success,
	})
	return
}
