package DomainLevel

import "time"

type GrpcInteraction interface {
	SendRequestGrpc([]byte) ([]byte, error)
	SayHi() string
}

type PacketChecker interface {
	Handle([]byte) (time.Duration, error)
	CheckSign([]byte, []byte, []byte) error
}

type GrpcDataManage interface {
	GenerateSignature(message []byte, key []byte) ([]byte, error)
	SayHI() string
}
