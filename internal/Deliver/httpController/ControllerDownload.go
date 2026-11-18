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

type DownloadController struct {
	DownloadNetwork
	Application.DownloadApplication
}

func GetNewDownload(downloadNetwork DownloadNetwork, downloadApplication Application.DownloadApplication) *DownloadController {
	return &DownloadController{DownloadNetwork: downloadNetwork, DownloadApplication: downloadApplication}
}

func (d DownloadController) Download() {
	name := mux.Vars(d.R)["name"]
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

	err := d.DownloadApplication.DownloadApplication(d.R.Context(), name, d.W)
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
