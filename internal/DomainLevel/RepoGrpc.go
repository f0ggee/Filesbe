package DomainLevel

type Requests interface {
	SetNewKeyRequest([]byte) ([]byte, error)
}
