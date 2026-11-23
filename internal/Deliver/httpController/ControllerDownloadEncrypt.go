package httpController

import (
	"Kaban/internal/Service/Application"
	"net/http"

	"github.com/gorilla/mux"
)

type DownloadEncryptController struct {
	Application.DownloadEncryptApplication
}

func (d DownloadEncryptController) DownloadEncrypt(W http.ResponseWriter, R *http.Request) {
	fileName := mux.Vars(R)["name"]
	if fileName == "" {
		SetAnswer(InputAnswerData{
			W:    W,
			Code: http.StatusCreated,
			Data: AnswerFileDownloadEncrypt{
				StatusOperation: Break,
				Error:           ErrorCantGetFileName.Error(),
			},
		})
		return
	}
	err := d.DownloadEncryptApplication.DownloadEncryptApplication(Application.DownloadEncryptIncomeData{
		Ctx:           R.Context(),
		EncryptedName: fileName,
	})
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    W,
			Code: http.StatusNotFound,
			Data: AnswerFileDownloadEncrypt{
				StatusOperation: NotStart,
				Error:           err.Error(),
			},
		})
		return
	}
	return
}
