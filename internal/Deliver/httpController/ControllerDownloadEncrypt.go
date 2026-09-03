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
type NewDownloadWithEncryptApplication struct {
	Application.NewDownloadEncrypt
}
type NewDownloadEncrypt struct {
	Net NetworkDownloadEncrypt
	NewDownloadWithEncryptApplication
	GetDataRequest func(r *http.Request) string
}

func GetNewDownloadEncrypt(net NetworkDownloadEncrypt, newDownloadWithEncryptApplication NewDownloadWithEncryptApplication) *NewDownloadEncrypt {
	return &NewDownloadEncrypt{Net: net, NewDownloadWithEncryptApplication: newDownloadWithEncryptApplication}
}

func (d NewDownloadEncrypt) DownloadWithEncrypt() {

	fileName := d.GetDataRequest(d.Net.R)
	if fileName == "" {
		SetAnswer(InputAnswerData{
			W:    d.Net.W,
			code: http.StatusCreated,
			data: AnswerFileDownloadEncrypt{
				StatusOperation: Break,
				Error:           ErrorCantGetFileName,
			},
		})
		return
	}
	err := d.DownloadEncrypt(Application.NewDownloadEncryptIncomingData{NewDownloadEncryptNetwork: Application.NewDownloadEncryptNetwork{d.Net.W},
		Ctx:          d.Net.R.Context(),
		EncryptedURl: fileName,
	})
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    d.Net.W,
			code: http.StatusNotFound,
			data: AnswerFileDownloadEncrypt{
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
