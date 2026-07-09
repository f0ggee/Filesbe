package DomainLevel

import "time"

type Swaping interface {
	SetSwapKeys() time.Duration
	SetSwapKeysFirst() time.Duration
	MakerRequests([]byte) time.Duration
}
