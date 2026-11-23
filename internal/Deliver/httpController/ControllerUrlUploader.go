package httpController

import (
	"Kaban/internal/DomainLevel"
	"log/slog"
	"net/http"
)

type BuildUrlController struct {
	Url func(r *http.Request) OutComingUrlData
}

func GetNewBuildUrl() *BuildUrlController {
	return &BuildUrlController{Url: UrlBuilder}
}

func (d *BuildUrlController) SetUrl(W http.ResponseWriter, R *http.Request) {
	urlData := d.Url(R)
	if urlData.NameFile == "" {
		slog.Error("UrlUploader name file empty", slog.Group("Request details", slog.String("URL", R.RequestURI)))
		SetAnswer(InputAnswerData{
			W: W,
			Data: AnswerUrlBuilder{
				StatusOperation: NotStart,
				ErrorMessage:    ErrorFileNameEmpty.Error(),
			},
		})
		return
	}
	switch {
	case urlData.FileType == "true":
		SetAnswer(InputAnswerData{
			W:    W,
			Code: http.StatusOK,
			Data: AnswerUrlBuilder{

				StatusOperation: Success,
				Url:             EncryptURLDownload + urlData.NameFile,
				ErrorMessage:    "",
			},
		})

		return

	case urlData.FileType == "false":
		SetAnswer(InputAnswerData{
			W:    W,
			Code: http.StatusOK,

			Data: AnswerUrlBuilder{
				StatusOperation: Success,
				Url:             UrlDownload,
			},
		})
		return
	}

	SetAnswer(InputAnswerData{
		W:    W,
		Code: http.StatusBadRequest,
		Data: AnswerUrlBuilder{
			ErrorMessage: DomainLevel.NonIdentifyError,
		},
	})
}

func UrlBuilder(r *http.Request) OutComingUrlData {
	name := r.URL.Query().Get(FileUrlName)
	boolParametric := r.URL.Query().Get(TypeFile)
	return OutComingUrlData{
		NameFile: name,
		FileType: boolParametric,
	}
}
