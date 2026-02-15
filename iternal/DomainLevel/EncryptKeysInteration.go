package DomainLevel

type EncryptionKeyInteration interface {
	ConvertDataToHash([]byte, []byte) []byte
	CheckSignIncomingKey([]byte, []byte, []byte) error
	DecryptIncomingKey([]byte, []byte, []byte) []byte
}
