package httpController

import (
	"Kaban/internal/DomainLevel"

	"Kaban/internal/Service/Application"
	"net/http"
)

type FileUploaderSessions struct {
	Rf       DomainLevel.AuthMaker
	Jwt      DomainLevel.AuthMaker
	Sessions Session
}
type UploadController struct {
	FileUploaderSessions
	Application.UploadApplication
}

func GetNewFileUploader(uploadNotEncryptSessions FileUploaderSessions, uploadApplication Application.UploadApplication) *UploadController {
	return &UploadController{FileUploaderSessions: uploadNotEncryptSessions, UploadApplication: uploadApplication}
}

func (d *UploadController) FileUploader(W http.ResponseWriter, R *http.Request) {

	returnedData := d.Sessions.GetSessionData(IncomingSessionData{Writer: W, Request: R})
	if returnedData.Error != nil {
		SetAnswer(InputAnswerData{
			W:    W,
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
			W:    W,
			Code: http.StatusUnauthorized,
			Data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				Error:           ErrorStrangeError.Error(),
			},
		})
		return
	}
	rfMaker, err := d.Rf.Make()
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    W,
			Code: http.StatusUnauthorized,
			Data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				Error:           ErrorStrangeError.Error(),
			},
		})
		return
	}

	jwtError := jwtMaker.IsTokenCorrect([]byte(returnedData.Jwt))
	rfError := rfMaker.IsTokenCorrect([]byte(returnedData.Rft))
	if rfError != nil && jwtError != nil {
		SetAnswer(InputAnswerData{
			W:    W,
			Code: http.StatusUnauthorized,
			Data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				Error:           ErrorAuthExpired.Error(),
			},
		})
		return
	}
	file, fileDetails, err := R.FormFile("File")
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    W,
			Code: http.StatusBadRequest,
			Data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				Error:           ErrorFile.Error(),
			},
		})
		return
	}

	fileName, err := d.UploadApplication.UploadApplication(Application.UploadIncomeData{
		File: file,
		Name: fileDetails.Filename,
		Size: fileDetails.Size,
		Ctx:  R.Context(),
	})
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    W,
			Code: http.StatusBadRequest,
			Data: AnswerUploaderFileNoEncrypt{
				StatusOperation: Break,
				Error:           err.Error(),
			},
		})
		return
	}
	SetAnswer(InputAnswerData{
		W:    W,
		Code: http.StatusCreated,
		Data: AnswerUploaderFileNoEncrypt{
			StatusOperation: Success,
			Params: ParametersUploading{
				IsEncrypt: false,
				FileName:  fileName,
			},
		},
	})
}
