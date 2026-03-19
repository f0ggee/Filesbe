package DomainLevel

import "time"

type SendingRequestGrpc interface {
	RequestingGettingNewKey([]byte) ([]byte, error)
	SayHi() string
}

type PacketChecker interface {
	CheckSign([]byte, []byte, []byte) error
	CheckTime(time.Time) error
}

type HandlingRequests interface {
	CheckingGettingNewKey([]byte) (time.Duration, error)
}
