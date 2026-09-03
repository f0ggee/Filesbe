package httpController

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type DownloadNetwork struct {
	W http.ResponseWriter
	R *http.Request
}

type DownloadApp struct {
	Download func(name string, ctx context.Context) error
}

type DownloadNew struct {
	DownloadNetwork
	DownloadApp
}

func GetNewDownloadWithNotEncrypt(networkDownloadNoEncrypt DownloadNetwork, newDownloadWithNotEncryptApplication DownloadApp) *DownloadNew {
	return &DownloadNew{DownloadNetwork: networkDownloadNoEncrypt, DownloadApp: newDownloadWithNotEncryptApplication}
}

func (d *DownloadNew) DownloadWithNotEncrypt() {
	name := d.getData(d.R)
	if name == "" {
		SetAnswer(InputAnswerData{
			W:    d.W,
			code: http.StatusBadRequest,
			data: AnswerDownloadNoEncrypt{
				StatusOperation: Break,
				Error:           ErrorCantGetFileName,
			},
		})
		return
	}

	err := d.Download(name, d.R.Context())
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    d.W,
			code: http.StatusBadRequest,
			data: AnswerDownloadNoEncrypt{
				StatusOperation: Break,
				Error:           err.Error(),
			},
		})
		return
	}
}
func (n *DownloadNew) getData(request *http.Request) string {
	vars := mux.Vars(request)
	name := vars["name"]
	return name
}
