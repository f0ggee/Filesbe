package RepourlBuilder

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"log/slog"
	"net/http"
)

type IncomingUrlData struct {
	Error           string
	StatusOperation string
	W               http.ResponseWriter
	Url             string
}

type OutComingData struct {
	NameFile string
	FileType string
}
type UrlBuilderAnswer interface {
	SetGoodAnswer(IncomingUrlData)
	SetBadAnswer(IncomingUrlData)
	GetUrlData(r *http.Request) OutComingData
}
type NewUrlBuilder struct{}

func GetNewUrlBuilder() *NewUrlBuilder {
	return &NewUrlBuilder{}
}

func (n NewUrlBuilder) SetGoodAnswer(data IncomingUrlData) {
	FilledAnswer := DomainLevel.AnswerUrlBuilder{
		StatusOperation: data.StatusOperation,
		Url:             data.Url,
	}
	data.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	data.W.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(data.W).Encode(FilledAnswer); err != nil {
		slog.Error("SetUrlBuilderGoodAnswer; error during encoding", "ERROR", err)
		return
	}
}

func (n NewUrlBuilder) SetBadAnswer(data IncomingUrlData) {
	FilledAnswer := DomainLevel.AnswerUrlBuilder{
		StatusOperation: data.StatusOperation,
		ErrorMessage:    data.Error,
	}
	data.W.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
	data.W.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(data.W).Encode(FilledAnswer); err != nil {
		slog.Error("SetUrlBuilderGoodAnswer; error during encoding", "ERROR", err)
		return
	}
}

func (n NewUrlBuilder) GetUrlData(r *http.Request) OutComingData {
	name := r.URL.Query().Get(DomainLevel.FileUrlName)
	boolParametric := r.URL.Query().Get(DomainLevel.TypeFile)
	return OutComingData{
		NameFile: name,
		FileType: boolParametric,
	}
}
