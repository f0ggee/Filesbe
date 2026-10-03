package DomainLevel

import (
	"io"
)

type Decoder interface {
	DecodeFlow(dst any, src io.Reader) error
	Decode(dst any, src []byte) error
}
type Encoder interface {
	EncodeFlow(dst any, src io.Writer) error
	Encode(dst any) ([]byte, error)
}
