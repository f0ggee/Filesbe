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
}

type NewDownload struct {
	DownloadNetwork
	Application.DownloadApplication
}

func GetNewDownload(downloadNetwork DownloadNetwork, downloadApplication Application.DownloadApplication) *NewDownload {
	return &NewDownload{DownloadNetwork: downloadNetwork, DownloadApplication: downloadApplication}
}

func (d *NewDownload) DownloadWithNotEncrypt() {
	name := d.getData(d.R)
	if name == "" {
		SetAnswer(InputAnswerData{
			W:    d.W,
			Code: http.StatusBadRequest,
			Data: AnswerDownloadNoEncrypt{
				StatusOperation: Break,
				Error:           ErrorCantGetFileName.Error(),
			},
		})
		return
	}

	err := d.Download(d.R.Context(), name)
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    d.W,
			Code: http.StatusBadRequest,
			Data: AnswerDownloadNoEncrypt{
				StatusOperation: Break,
				Error:           err.Error(),
			},
		})
		return
	}
}
func (n *NewDownload) getData(request *http.Request) string {
	vars := mux.Vars(request)
	name := vars["name"]
	return name
}
