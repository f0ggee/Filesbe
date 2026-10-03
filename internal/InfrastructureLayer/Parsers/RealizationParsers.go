package Parsers

import (
	"Kaban/internal/DomainLevel"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
)

type Parsing struct {
}

func GetNewParsing() Parsing {
	return Parsing{}
}

func (p Parsing) EncodeFlow(a any, writer io.Writer) error {

	if err := json.NewEncoder(writer).Encode(&a); err != nil {
		slog.Error("GetIncomingData: the error to parse data", "ERROR", err)
		return errors.New(DomainLevel.ErrorParseInfo)
	}
	return nil
}

func (p Parsing) Encode(a any) ([]byte, error) {

	marshal, err := json.Marshal(&a)
	if err != nil {
		slog.Error("Encode; error to parse data", "ERROR", err)
		return nil, errors.New(DomainLevel.ErrorParseInfo)
	}
	return marshal, nil
}
func (p Parsing) DecodeFlow(a any, request io.Reader) error {

	if err := json.NewDecoder(request).Decode(&a); err != nil {
		slog.Error("GetIncomingData: the error to parse data", "ERROR", err)
		return errors.New(DomainLevel.ErrorParseInfo)
	}
	return nil
}

func (p Parsing) Decode(a any, bytes []byte) error {

	err := json.Unmarshal(bytes, &a)
	if err != nil {
		slog.Error("Decoder; error to parse data", "ERROR", err)
		return errors.New(DomainLevel.ErrorParseInfo)
	}
	return nil
}
