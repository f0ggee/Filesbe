package httpController

import (
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"Kaban/internal/InfrastructureLayer/RepoSessionHandle"
	"Kaban/internal/Service/Application"
	"errors"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/gorilla/mux"
)

type NewFileUploaderEncryptNetwork struct {
	W http.ResponseWriter
	R *http.Request
}

type NewFileUploaderEncryptSession struct {
	ReadSession RepoSessionHandle.Session
	AuthCheck   AuthTokensManage.AuthCheck
}
type NewFileUploaderEncryptDetails struct {
	Rout *mux.Router
}
type NewUploaderEncrypt struct {
	NewFileUploaderEncryptNetwork
	NewFileUploaderEncryptSession
	NewFileUploaderEncryptDetails
	UploadEncrypt func(data Application.IncomeData) (string, error)
	UrlUploadData func(r *mux.Router, fileName string) (string, error)
}

func GetNewFileUploaderEncrypt(fileUploaderEncryptNetwork NewFileUploaderEncryptNetwork, newFileUploaderEncryptSession NewFileUploaderEncryptSession, newFileUploaderEncryptDetails NewFileUploaderEncryptDetails) *NewUploaderEncrypt {
	return &NewUploaderEncrypt{NewFileUploaderEncryptNetwork: fileUploaderEncryptNetwork, NewFileUploaderEncryptSession: newFileUploaderEncryptSession, NewFileUploaderEncryptDetails: newFileUploaderEncryptDetails}
}

func (S *NewUploaderEncrypt) FileUploaderEncrypt() {

	returnedSession := S.ReadSession.GetSessionData(RepoSessionHandle.IncomingSessionData{
		Writer:  S.W,
		Request: S.R,
	})
	if returnedSession.Error != nil {

		SetAnswer(InputAnswerData{
			W:    S.W,
			code: http.StatusBadRequest,
			data: AnswerUploadEncrypt{
				StatusOperation: Break,
				Error:           returnedSession.Error.Error(),
			},
		})
		return
	}
	Data := S.AuthCheck.CheckUserAuth(AuthTokensManage.UserAuthCheckIncomingData{
		Jwt: returnedSession.Jwt,
		Rft: returnedSession.Rft,
	})
	if Data.Err != nil {
		SetAnswer(InputAnswerData{
			W:    S.W,
			code: http.StatusUnauthorized,
			data: AnswerUploadEncrypt{
				StatusOperation: Break,
				Error:           Data.Err.Error(),
			},
		})
		return
	}
	file, sizeAndName, err := S.getFile()
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    S.W,
			code: http.StatusBadRequest,
			data: errors.New(ErrorFile),
		})
		return
	}
	fileName, err := S.UploadEncrypt(Application.IncomeData{
		File: file,
		Name: sizeAndName.Filename,
		Size: sizeAndName.Size,
		Ctx:  S.R.Context(),
	})
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    S.W,
			code: http.StatusBadRequest,
			data: AnswerUploadEncrypt{
				StatusOperation: NotStart,
				Error:           err.Error(),
			},
		})
		return
	}

	urlPath, err := S.UrlUploadData(S.Rout, fileName)
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    S.W,
			code: http.StatusBadRequest,
			data: AnswerUploadEncrypt{
				StatusOperation: NotStart,
				Error:           err.Error(),
			},
		})
		return
	}
	SetAnswer(InputAnswerData{
		W:    S.W,
		code: http.StatusOK,
		data: AnswerUploadEncrypt{
			StatusOperation: Success,
			UrlToRedirect:   urlPath,
		},
	})
	return
}

func (S *NewUploaderEncrypt) getFile() (multipart.File, *multipart.FileHeader, error) {
	file, sizeAndName, err := S.R.FormFile("file")
	if err != nil {
		slog.Error("UploadEncrypt; error to get a file", "ERROR", err)
		return nil, nil, err
	}
	return file, sizeAndName, nil
}

func getUploadEncryptData(r *mux.Router, fileName string) (string, error) {
	url, err := r.Get("fileName").URL("name", fileName, "bool", "true")
	if err != nil {
		slog.Error("UrlBuilderUploadEncrypt; error to get a file name from the url", "ERROR", err)
		return "", errors.New(ErrorCantGetFileName)
	}
	return url.Path, nil
}
