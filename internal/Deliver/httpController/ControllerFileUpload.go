package httpController

import (
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"Kaban/internal/InfrastructureLayer/RepoSessionHandle"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type FileUploaderNet struct {
	W http.ResponseWriter
	R *http.Request
}

type FileUploaderSessions struct {
	Session RepoSessionHandle.Session
	Auth    AuthTokensManage.AuthCheck
}
type NewFileUploader struct {
	FileUploaderNet
	FileUploaderSessions
	FileUploader func(r *http.Request) (string, error)
	UrlData      func(r *mux.Router, fileName string) (string, error)
}

func GetNewFileUploader(fileUploaderNoEncryptNet FileUploaderNet, uploadNotEncryptSessions FileUploaderSessions) *NewFileUploader {
	return &NewFileUploader{FileUploaderNet: fileUploaderNoEncryptNet, FileUploaderSessions: uploadNotEncryptSessions}
}

func (d *NewFileUploader) FileUploaderNoEncrypt(router *mux.Router) {
	returnedData := d.Session.GetSessionData(RepoSessionHandle.IncomingSessionData{Writer: d.W, Request: d.R})
	if returnedData.Error != nil {
		SetAnswer(InputAnswerData{
			W:    d.W,
			code: http.StatusUnauthorized,
			data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				Error:           returnedData.Error.Error(),
			},
		})
		return
	}
	outData := d.Auth.CheckUserAuth(AuthTokensManage.UserAuthCheckIncomingData{
		Jwt: returnedData.Jwt,
		Rft: returnedData.Rft,
	})
	if outData.Err != nil {
		SetAnswer(InputAnswerData{
			W:    d.W,
			code: http.StatusUnauthorized,
			data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				Error:           outData.Err.Error(),
			},
		})
		return
	}
	if outData.IsNewJwtCreated {
		d.Session.SetNewSession(RepoSessionHandle.IncomingSessionData{Jwt: outData.NewJwt})
	}

	fileName, err := d.FileUploader(d.R)
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    d.W,
			code: http.StatusBadRequest,
			data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				Error:           err.Error(),
			},
		})
		return
	}

	urlPath, err := d.UrlData(router, fileName)
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    d.W,
			code: http.StatusBadRequest,
			data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				Error:           err.Error(),
			},
		})
		return
	}
	SetAnswer(InputAnswerData{
		W:    d.W,
		code: http.StatusCreated,
		data: AnswerUploaderFileNoEncrypt{
			StatusOperation: Success,
			UrlToRedirect:   urlPath,
		},
	})
	return
}

func getUploadData(r *mux.Router, fileName string) (string, error) {
	url, err := r.Get("fileName").URL("name", fileName, "bool", "true")
	if err != nil {
		slog.Error("UrlBuilderUploadEncrypt; error to get a file name from the url", "ERROR", err)
		return "", errors.New(ErrorCantGetFileName)
	}
	return url.Path, nil
}
