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
type NewDownloadEncrypt struct {
	Net NetworkDownloadEncrypt
	Application.DownloadEncryptApplication
	GetDataRequest func(r *http.Request) string
}

func (d NewDownloadEncrypt) DownloadWithEncrypt() {

	fileName := d.GetDataRequest(d.Net.R)
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
	err := d.DownloadEncrypt(Application.NewDownloadEncryptIncomingData{NewDownloadEncryptNetwork: Application.NewDownloadEncryptNetwork{d.Net.W},
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

func GetDataRequest(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["name"]
}
