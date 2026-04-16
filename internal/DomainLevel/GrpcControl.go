package DomainLevel

import "time"

type SendRequestGrpc interface {
	RequestingGettingNewKey([]byte) ([]byte, error)
	SayHi() string
}

type PacketChecker interface {
	CheckTime(time.Time) error
}

type HandlingRequests interface {
	CheckingGettingNewKey([]byte) (time.Duration, error)
}

type GrpcTest interface {
	EncryptByWrongKey([]byte) ([]byte, error)
}
