package DomainLevel

import (
	"io"
)

type Decoder interface {
	DecodeFlow(any, io.Reader) error
	Decode(any, []byte) error
}
type Encoder interface {
	EncodeFlow(any, io.Writer) error
	Encode(any) ([]byte, error)
}
