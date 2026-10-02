package httpController

import (
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"Kaban/internal/InfrastructureLayer/RepoSession"
	"Kaban/internal/Service/Application"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/gorilla/mux"
)

type FileUploaderNet struct {
	W http.ResponseWriter
	R *http.Request
}

type FileUploaderSessions struct {
	Session RepoSession.Session
	Auth    AuthTokensManage.AuthCheck
}
type NewFileUploader struct {
	FileUploaderNet
	FileUploaderSessions
	Application.UploadApplication
	UrlData func(r *mux.Router, fileName string) (string, error)
}

func GetNewFileUploader(fileUploaderNoEncryptNet FileUploaderNet, uploadNotEncryptSessions FileUploaderSessions) *NewFileUploader {
	return &NewFileUploader{FileUploaderNet: fileUploaderNoEncryptNet, FileUploaderSessions: uploadNotEncryptSessions}
}

func (d *NewFileUploader) FileUploaderNoEncrypt(router *mux.Router) {
	returnedData := d.Session.GetSessionData(RepoSession.IncomingSessionData{Writer: d.W, Request: d.R})
	if returnedData.Error != nil {
		SetAnswer(InputAnswerData{
			W:    d.W,
			Code: http.StatusUnauthorized,
			Data: AnswerUploaderFileNoEncrypt{
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
			Code: http.StatusUnauthorized,
			Data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				Error:           outData.Err.Error(),
			},
		})
		return
	}
	if outData.IsNewJwtCreated {
		d.Session.SetNewSession(RepoSession.IncomingSessionData{Jwt: outData.NewJwt})
	}
	file, fileDetails, err := d.getFileData()
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    d.W,
			Code: http.StatusBadRequest,
			Data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				Error:           ErrorCantGetFileName.Error(),
			},
		})
		return
	}
	fileName, err := d.Upload(Application.FileUploaderIncomeData{
		File: file,
		Name: fileDetails.Filename,
		Size: fileDetails.Size,
		Ctx:  d.R.Context(),
	})
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    d.W,
			Code: http.StatusBadRequest,
			Data: AnswerUploaderFileNoEncrypt{
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
			Code: http.StatusBadRequest,
			Data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				Error:           err.Error(),
			},
		})
		return
	}
	SetAnswer(InputAnswerData{
		W:    d.W,
		Code: http.StatusCreated,
		Data: AnswerUploaderFileNoEncrypt{
			StatusOperation: Success,
			UrlToRedirect:   urlPath,
		},
	})
	return
}

func (d *NewFileUploader) getFileData() (multipart.File, *multipart.FileHeader, error) {
	file, fileDetails, err := d.R.FormFile("File")
	if err != nil {
		slog.Error("Upload; error to get a File", "ERROR", err)
		return nil, nil, ErrorCantGetFileName
	}
	return file, fileDetails, nil
}

func getUploadData(r *mux.Router, fileName string) (string, error) {
	url, err := r.Get("fileName").URL("name", fileName, "bool", "true")
	if err != nil {
		slog.Error("UrlBuilderUploadEncrypt; error to get a file name from the url", "ERROR", err)
		return "", ErrorCantGetFileName
	}
	return url.Path, nil
}
