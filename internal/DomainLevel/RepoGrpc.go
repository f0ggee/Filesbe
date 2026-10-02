package DomainLevel

type Requests interface {
	SetNewKeyRequest([]byte) ([]byte, error)
}

type MakerKeyReqeust interface {
	SetAdditionalData([]byte) MakerKeyReqeust
	Make() (Requests, error)
}
