package DomainLevel

type Requests interface {
	SetEncrypterKeyRequest([]byte) ([]byte, error)
	SetMakerRequestEncrypterKey([]byte) ([]byte, error)
}
