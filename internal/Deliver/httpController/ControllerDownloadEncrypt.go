package httpController

import (
	"Kaban/internal/Service/Application"
	"net/http"

	"github.com/gorilla/mux"
)

type NetworkDownloadEncrypt struct {
	W http.ResponseWriter
	R *http.Request
}
type DownloadEncryptController struct {
	Net NetworkDownloadEncrypt
	Application.DownloadEncryptApplication
}

func (d DownloadEncryptController) DownloadEncrypt() {
	fileName := mux.Vars(d.Net.R)["name"]
	if fileName == "" {
		SetAnswer(InputAnswerData{
			W:    d.Net.W,
			Code: http.StatusCreated,
			Data: AnswerFileDownloadEncrypt{
				StatusOperation: Break,
				Error:           ErrorCantGetFileName.Error(),
			},
		})
		return
	}
	err := d.DownloadEncryptApplication.DownloadEncryptApplication(Application.DownloadEncryptIncomeData{
		Ctx:           d.Net.R.Context(),
		EncryptedName: fileName,
	})
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    d.Net.W,
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
