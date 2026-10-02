package DomainLevel

type Requests interface {
	SetNewKeyRequest([]byte) ([]byte, error)
}

type MakerKeyRequest interface {
	SetAdditionalData([]byte) MakerKeyRequest
	Make() (Requests, error)
}
