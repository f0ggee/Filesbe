package DomainLevel

import "context"

type Requests interface {
	SetNewKeyRequest([]byte) ([]byte, error)
}

type MakerKeyRequest interface {
	SetAdditionalData([]byte) MakerKeyRequest
	Make(context.Context) (Requests, error)
}
