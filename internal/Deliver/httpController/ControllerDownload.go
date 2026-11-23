package httpController

import (
	"Kaban/internal/Service/Application"
	"net/http"

	"github.com/gorilla/mux"
)

type DownloadController struct {
	Application.DownloadApplication
}

func GetNewDownload(downloadApplication Application.DownloadApplication) *DownloadController {
	return &DownloadController{DownloadApplication: downloadApplication}
}

func (d DownloadController) Download(W http.ResponseWriter, R *http.Request) {
	name := mux.Vars(R)["name"]
	if name == "" {
		SetAnswer(InputAnswerData{
			W:    W,
			Code: http.StatusBadRequest,
			Data: AnswerDownloadNoEncrypt{
				StatusOperation: Break,
				Error:           ErrorCantGetFileName.Error(),
			},
		})
		return
	}

	err := d.DownloadApplication.DownloadApplication(R.Context(), name, W)
	if err != nil {
		SetAnswer(InputAnswerData{
			W:    W,
			Code: http.StatusBadRequest,
			Data: AnswerDownloadNoEncrypt{
				StatusOperation: Break,
				Error:           err.Error(),
			},
		})
		return
	}
}
