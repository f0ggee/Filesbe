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
type NewFileUploaderEncryptSession struct {
	ReadSession RepoSessionHandle.Session
	AuthCheck   AuthTokensManage.AuthCheck
}
type NewFileUploaderEncryptDetails struct {
	Answers RepofileUploaderEncryptRepo.AnswersUploadEncrypt
	Build   RepofileUploaderEncryptRepo.UrlUploadEncrypt
	Rout    *mux.Router
}
type NewFileUploaderEncrypt struct {
	FileUploaderEncryptNetwork
	NewFileUploaderEncryptSession
	NewFileUploaderEncryptDetails
	NewFileUploaderEncryptApplication
}

func GetNewFileUploaderEncrypt(fileUploaderEncryptNetwork FileUploaderEncryptNetwork, newFileUploaderEncryptSession NewFileUploaderEncryptSession, newFileUploaderEncryptDetails NewFileUploaderEncryptDetails, newFileUploaderEncryptApplication NewFileUploaderEncryptApplication) *NewFileUploaderEncrypt {
	return &NewFileUploaderEncrypt{FileUploaderEncryptNetwork: fileUploaderEncryptNetwork, NewFileUploaderEncryptSession: newFileUploaderEncryptSession, NewFileUploaderEncryptDetails: newFileUploaderEncryptDetails, NewFileUploaderEncryptApplication: newFileUploaderEncryptApplication}
}

func (S *NewFileUploaderEncrypt) FileUploaderEncrypt() {

	returnedSession := S.ReadSession.GetSessionData(RepoSessionHandle.IncomingSessionData{
		Writer:  S.W,
		Request: S.R,
	})
	if returnedSession.Error != nil {
		S.Answers.SetBadAnswers(RepofileUploaderEncryptRepo.IncomingDataAnswer{
			W:               S.W,
			Error:           returnedSession.Error.Error(),
			StatusOperation: DomainLevel.Break,
		})
		return
	}
	Data := S.AuthCheck.CheckUserAuth(AuthTokensManage.UserAuthCheckIncomingData{
		Jwt: returnedSession.Jwt,
		Rft: returnedSession.Rft,
	})
	if Data.Err != nil {
		S.Answers.SetGoodAnswers(RepofileUploaderEncryptRepo.IncomingDataAnswer{
			W:               S.W,
			Error:           Data.Err.Error(),
			StatusOperation: DomainLevel.Break,
		})
		return
	}

	fileName, err := S.NewUploadEncrypt.UploadEncrypt(S.R)
	if err != nil {
		//TODO add handling the error

		return
	}

	urlPath, err := S.Build.UrlBuilder(S.Rout, fileName)
	if err != nil {
		S.Answers.SetBadAnswers(RepofileUploaderEncryptRepo.IncomingDataAnswer{
			W:               S.W,
			Error:           err.Error(),
			StatusOperation: DomainLevel.NotStart,
		})
		return
	}

	S.Answers.SetGoodAnswers(RepofileUploaderEncryptRepo.IncomingDataAnswer{
		W:               S.W,
		UrlToRedirect:   urlPath,
		StatusOperation: DomainLevel.Success,
	})
	return
}
