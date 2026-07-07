package DomainLevel

import "time"

type SendRequestGrpc interface {
	RequestingGettingNewKey([]byte) ([]byte, error)
}

type HandlingRequests interface {
	CheckingGettingNewKey([]byte) (time.Duration, error)
}
