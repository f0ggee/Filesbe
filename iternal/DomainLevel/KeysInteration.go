package DomainLevel

type KeyInteraction interface {
	ConvertDataToHash([]byte, []byte) []byte
	CheckSignIncomingKey([]byte, []byte, []byte) error
	DecryptIncomingKey([]byte, []byte, []byte) []byte
	JsonConverter(any) ([]byte, error)
}

type EncryptionKey interface {
	EncryptAes([]byte, []byte) ([]byte, error)
}
