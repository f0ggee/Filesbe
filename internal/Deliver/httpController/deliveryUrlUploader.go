package httpController

import (
	"Kaban/internal/DomainLevel"
	"log/slog"
	"net/http"
)

type UrlNetwork struct {
	W http.ResponseWriter
	r *http.Request
}
type NewBuildUrl struct {
	Url func(r *http.Request) OutComingUrlData
	Net UrlNetwork
}

func GetNewBuildUrl(net UrlNetwork) *NewBuildUrl {
	return &NewBuildUrl{Net: net}
}

func (d *NewBuildUrl) SetUrl() {
	urlData := d.Url(d.Net.r)
	if urlData.NameFile == "" {
		slog.Error("UrlUploader name file empty", slog.Group("Request details", slog.String("URL", d.Net.r.RequestURI)))
		SetAnswer(InputAnswerData{
			W: d.Net.W,
			data: AnswerUrlBuilder{
				StatusOperation: NotStart,
				ErrorMessage:    ErrorFileNameEmpty,
			},
		})
		return
	}
	switch {
	case urlData.FileType == "true":
		SetAnswer(InputAnswerData{
			W:    d.Net.W,
			code: http.StatusOK,
			data: AnswerUrlBuilder{

				StatusOperation: Success,
				Url:             EncryptURLDownload + urlData.NameFile,
				ErrorMessage:    "",
			},
		})

		return

	case urlData.FileType == "false":
		SetAnswer(InputAnswerData{
			W:    d.Net.W,
			code: http.StatusOK,

			data: AnswerUrlBuilder{
				StatusOperation: Success,
				Url:             UrlDownload,
			},
		})
		return
	}

	SetAnswer(InputAnswerData{
		W:    d.Net.W,
		code: http.StatusBadRequest,
		data: AnswerUrlBuilder{
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
