package RepoParsers

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

type Decode interface {
	JsonDecode(any, io.Reader) error
	JsonDecodeMarshall(any, []byte) error
}
type Encode interface {
	JsonEncode(any, http.ResponseWriter) error
	JsonEncodeMarshall(any) ([]byte, error)
}

type Parsing struct {
}

func GetNewParsing() *Parsing {
	return &Parsing{}
}

func (p Parsing) JsonEncode(a any, writer http.ResponseWriter) error {

	if err := json.NewEncoder(writer).Encode(&a); err != nil {
		slog.Error("GetIncomingData: the error to parse data", "ERROR", err)
		return errors.New(DomainLevel.ErrorParseInfo)
	}
	return nil
}

func (p Parsing) JsonEncodeMarshall(a any) ([]byte, error) {

	marshal, err := json.Marshal(&a)
	if err != nil {
		slog.Error("JsonEncodeMarshall; error to parse data", "ERROR", err)
		return nil, errors.New(DomainLevel.ErrorParseInfo)
	}
	return marshal, nil
}
func (p Parsing) JsonDecode(a any, request io.Reader) error {

	if err := json.NewDecoder(request).Decode(&a); err != nil {
		slog.Error("GetIncomingData: the error to parse data", "ERROR", err)
		return errors.New(DomainLevel.ErrorParseInfo)
	}
	return nil
}

func (p Parsing) JsonDecodeMarshall(a any, bytes []byte) error {

	err := json.Unmarshal(bytes, &a)
	if err != nil {
		slog.Error("JsonDecodeMarshall; error to parse data", "ERROR", err)
		return errors.New(DomainLevel.ErrorParseInfo)
	}
	return nil
}
