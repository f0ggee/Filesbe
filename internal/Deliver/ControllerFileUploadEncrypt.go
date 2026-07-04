package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage/AuthChecking"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoSessionHandle"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepofileUploaderEncryptRepo"
	"net/http"

	"github.com/gorilla/mux"
)

type FileUploaderEncryptNetwork struct {
	W http.ResponseWriter
	R *http.Request
}

type FileUploaderEncryptRouting struct {
	Rout *mux.Router
}
type Session struct {
	ReadSession RepoSessionHandle.SessionConnect
	AuthCheck   AuthChecking.NewAuthChecker
}
type AnswerUploadEncrypt struct {
	S RepofileUploaderEncryptRepo.SetNewUploadingRepo
}
type NewFileUploaderEncrypt struct {
	Net    FileUploaderNoEncryptNet
	Router FileUploaderEncryptRouting
	Sess   Session
	Answe  AnswerUploadEncrypt
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
	Data := S.Sess.AuthCheck.CheckAuthTokens(DomainLevel.AuthCheckIncomingData{
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

	filName, err := s.UploadEncrypt(r)
	if err != nil {
		//TODO add handling the error

		return
	}

	urlPath, err := S.Answe.S.UrlBuilder(S.Router.Rout, filName)
	if err != nil {
		S.Answe.S.SetBadAnswers(RepofileUploaderEncryptRepo.IncomingDataAnswer{
			W:               S.Net.w,
			Error:           err,
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
