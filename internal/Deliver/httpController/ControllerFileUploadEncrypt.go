package httpController

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Service/Application"
	"net/http"
)

type NewFileUploaderEncryptSession struct {
	ReadSession Session
	Jwt         DomainLevel.AuthMaker
	Rf          DomainLevel.AuthMaker
}

type UploadEncryptController struct {
	NewFileUploaderEncryptSession
	Application.UploaderEncryptApplication
}

func (S *UploadEncryptController) FileUploaderEncrypt(W http.ResponseWriter, R *http.Request) {
	returnedSession := S.ReadSession.GetSessionData(IncomingSessionData{
		Writer:  W,
		Request: R,
	})
	if returnedSession.Error != nil {
		SetAnswer(InputAnswerData{
			W:    W,
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
			W:    W,
			Code: http.StatusUnauthorized,
			Data: AnswerUploadEncrypt{
				StatusOperation: Break,
				Error:           ErrorStrangeError.Error(),
			},
		})
		return
	}
	rfMaker, err := S.Rf.Make()
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    W,
			Code: http.StatusUnauthorized,
			Data: AnswerUploadEncrypt{
				StatusOperation: Break,
				Error:           ErrorStrangeError.Error(),
			},
		})
		return
	}
	jwtError := jwtMaker.IsTokenCorrect([]byte(returnedSession.Jwt))
	rfError := rfMaker.IsTokenCorrect([]byte(returnedSession.Rft))
	if rfError != nil && jwtError != nil {
		SetAnswer(InputAnswerData{
			W:    W,
			Code: http.StatusUnauthorized,
			Data: AnswerUploadEncrypt{
				StatusOperation: Break,
				Error:           ErrorAuthExpired.Error(),
			},
		})
		return
	}
	file, sizeAndName, err := R.FormFile("File")
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    W,
			Code: http.StatusBadRequest,
			Data: AnswerUploadEncrypt{
				StatusOperation: Break,
				Error:           ErrorFile.Error(),
			},
		})
		return
	}
	fileName, err := S.UploadEncryptApplication(Application.UploadEncryptIncomeData{
		File: file,
		Name: []byte(sizeAndName.Filename),
		Size: sizeAndName.Size,
		Ctx:  R.Context(),
	})
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    W,
			Code: http.StatusBadRequest,
			Data: AnswerUploadEncrypt{
				StatusOperation: NotStart,
				Error:           err.Error(),
			},
		})
		return
	}

	SetAnswer(InputAnswerData{
		W:    W,
		Code: http.StatusOK,
		Data: AnswerUploadEncrypt{
			StatusOperation: Success,
			Params: ParametersUploading{
				IsEncrypt: true,
				FileName:  fileName,
			},
		},
	})
	return
}
