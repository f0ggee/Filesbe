package httpController

import (
	"Kaban/internal/Service/Application"
	"net/http"

	"github.com/gorilla/mux"
)

type DownloadNetwork struct {
	W http.ResponseWriter
	R *http.Request
}

type DownloadApp struct {
	Application.DownloadApplication
}

type DownloadNew struct {
	DownloadNetwork
	DownloadApp
}

func GetDownloadNew(downloadNetwork DownloadNetwork, downloadApp DownloadApp) *DownloadNew {
	return &DownloadNew{DownloadNetwork: downloadNetwork, DownloadApp: downloadApp}
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

	err := d.Download(d.R.Context(), name)
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
