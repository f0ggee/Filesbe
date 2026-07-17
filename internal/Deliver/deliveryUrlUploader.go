package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepourlBuilder"
	"log/slog"
	"net/http"
)

type UrlSettings struct {
	Url RepourlBuilder.UrlBuilderAnswer
}

type UrlNetwork struct {
	W http.ResponseWriter
	r *http.Request
}
type NewBuildUrl struct {
	S   UrlSettings
	Net UrlNetwork
}

func GetNewBuildUrl(s UrlSettings, net UrlNetwork) *NewBuildUrl {
	return &NewBuildUrl{S: s, Net: net}
}

func (d *NewBuildUrl) SetUrl() {
	urlData := d.S.Url.GetUrlData(d.Net.r)
	if urlData.NameFile == "" {
		slog.Error("UrlUploader name file empty", slog.Group("Request details", slog.String("URL", d.Net.r.RequestURI)))
		d.S.Url.SetBadAnswer(RepourlBuilder.IncomingUrlData{
			Error:           DomainLevel.ErrorFileNameEmpty,
			StatusOperation: DomainLevel.NotStart,
			W:               d.Net.W,
		})
		return
	}
	switch {
	case urlData.FileType == "true":
		url := DomainLevel.DomainName + "d2/" + urlData.NameFile

		d.S.Url.SetGoodAnswer(RepourlBuilder.IncomingUrlData{
			StatusOperation: DomainLevel.Success,
			W:               d.Net.W,
			Url:             url,
		})
		return

	case urlData.FileType == "false":
		d.S.Url.SetGoodAnswer(RepourlBuilder.IncomingUrlData{
			StatusOperation: DomainLevel.Success,
			W:               d.Net.W,
			Url:             DomainLevel.DomainName + "d/" + urlData.NameFile,
		})
		return
	}
}
