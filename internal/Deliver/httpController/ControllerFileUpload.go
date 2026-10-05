package httpController

import (
	"Kaban/internal/DomainLevel"

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
	Rf       DomainLevel.AuthMaker
	Jwt      DomainLevel.AuthMaker
	Sessions Session
}
type UploadController struct {
	FileUploaderNet
	FileUploaderSessions
	Application.UploadApplication
	UrlData func(r *mux.Router, fileName string) (string, error)
}

func GetNewFileUploader(fileUploaderNoEncryptNet FileUploaderNet, uploadNotEncryptSessions FileUploaderSessions) *UploadController {
	return &UploadController{FileUploaderNet: fileUploaderNoEncryptNet, FileUploaderSessions: uploadNotEncryptSessions}
}

func (d *UploadController) FileUploader(router *mux.Router) {

	returnedData := d.Sessions.GetSessionData(IncomingSessionData{Writer: d.W, Request: d.R})
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

	jwtMaker, err := d.Jwt.Make()
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    d.W,
			Code: http.StatusUnauthorized,
			Data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				UrlToRedirect:   "",
				Error:           ErrorStrangeError.Error(),
			},
		})
		return
	}
	rfMaker, err := d.Rf.Make()
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    d.W,
			Code: http.StatusUnauthorized,
			Data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				UrlToRedirect:   "",
				Error:           ErrorStrangeError.Error(),
			},
		})
		return
	}

	jwtError := jwtMaker.IsTokenCorrect([]byte(returnedData.Jwt))
	rfError := rfMaker.IsTokenCorrect([]byte(returnedData.Rft))
	if rfError != nil && jwtError != nil {
		SetAnswer(InputAnswerData{
			W:    d.W,
			Code: http.StatusUnauthorized,
			Data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				UrlToRedirect:   "/login",
				Error:           ErrorAuthExpired.Error(),
			},
		})

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

func (d *UploadController) getFileData() (multipart.File, *multipart.FileHeader, error) {
	file, fileDetails, err := d.R.FormFile("File")
	if err != nil {
		slog.Error("Upload; error to get a File", "ERROR", err)
		return nil, nil, ErrorCantGetFileName
	}
	return file, fileDetails, nil
}

func GetUploadData(r *mux.Router, fileName string) (string, error) {
	url, err := r.Get("fileName").URL("name", fileName, "bool", "true")
	if err != nil {
		slog.Error("UrlBuilderUploadEncrypt; error to get a file name from the url", "ERROR", err)
		return "", ErrorCantGetFileName
	}
	return url.Path, nil
}
