package httpController

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Service/Application"
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
	ReadSession Session
	Jwt         DomainLevel.AuthMaker
	Rf          DomainLevel.AuthMaker
}

type UploadEncryptController struct {
	NewFileUploaderEncryptNetwork
	NewFileUploaderEncryptSession
	Application.UploaderEncryptApplication
	Rout *mux.Router

	UrlUploadData func(r *mux.Router, fileName string) (string, error)
}

func (S *UploadEncryptController) FileUploaderEncrypt() {

	returnedSession := S.ReadSession.GetSessionData(IncomingSessionData{
		Writer:  S.W,
		Request: S.R,
	})
	if returnedSession.Error != nil {

		SetAnswer(InputAnswerData{
			W:    S.W,
			Code: http.StatusBadRequest,
			Data: AnswerUploadEncrypt{
				StatusOperation: Break,
				Error:           returnedSession.Error.Error(),
			},
		})
		return
	}
	jwtMaker, err := S.Jwt.Make()
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    S.W,
			Code: http.StatusUnauthorized,
			Data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				UrlToRedirect:   "",
				Error:           ErrorStrangeError.Error(),
			},
		})
		return
	}
	rfMaker, err := S.Rf.Make()
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    S.W,
			Code: http.StatusUnauthorized,
			Data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				UrlToRedirect:   "",
				Error:           ErrorStrangeError.Error(),
			},
		})
		return
	}
	jwtError := jwtMaker.IsTokenCorrect([]byte(returnedSession.Jwt))
	rfError := rfMaker.IsTokenCorrect([]byte(returnedSession.Rft))
	if rfError != nil && jwtError != nil {
		SetAnswer(InputAnswerData{
			W:    S.W,
			Code: http.StatusUnauthorized,
			Data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				UrlToRedirect:   "/login",
				Error:           ErrorAuthExpired.Error(),
			},
		})

	}
	file, sizeAndName, err := S.getFile()
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    S.W,
			Code: http.StatusBadRequest,
			Data: ErrorFile,
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
			Code: http.StatusBadRequest,
			Data: AnswerUploadEncrypt{
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
			Code: http.StatusBadRequest,
			Data: AnswerUploadEncrypt{
				StatusOperation: NotStart,
				Error:           err.Error(),
			},
		})
		return
	}
	SetAnswer(InputAnswerData{
		W:    S.W,
		Code: http.StatusOK,
		Data: AnswerUploadEncrypt{
			StatusOperation: Success,
			UrlToRedirect:   urlPath,
		},
	})
	return
}

func (S *UploadEncryptController) getFile() (multipart.File, *multipart.FileHeader, error) {
	file, sizeAndName, err := S.R.FormFile("file")
	if err != nil {
		slog.Error("UploadEncrypt; error to get a file", "ERROR", err)
		return nil, nil, err
	}
	return file, sizeAndName, nil
}

func GetUploadEncryptData(r *mux.Router, fileName string) (string, error) {
	url, err := r.Get("fileName").URL("name", fileName, "bool", "true")
	if err != nil {
		slog.Error("UrlBuilderUploadEncrypt; error to get a file name from the url", "ERROR", err)
		return "", ErrorCantGetFileName
	}
	return url.Path, nil
}
