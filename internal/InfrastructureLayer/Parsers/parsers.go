package Parsers

import (
	"Kaban/internal/DomainLevel"
	"crypto/rand"
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
		slog.Error("Parser; error to parse data", "ERROR", err)
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

type ParseError struct {
}

func (p ParseError) DecodeFlow(dst any, src io.Reader) error {
	return errors.New(rand.Text())
}

func (p ParseError) Decode(dst any, src []byte) error {
	return errors.New(rand.Text())
}

func (p ParseError) EncodeFlow(dst any, src io.Writer) error {
	return errors.New(rand.Text())
}

func (p ParseError) Encode(dst any) ([]byte, error) {
	return nil, errors.New(rand.Text())
}

func GetParseError() *ParseError {
	return &ParseError{}
}
