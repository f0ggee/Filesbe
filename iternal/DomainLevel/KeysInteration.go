package DomainLevel

import "crypto/rsa"

type KeyInteraction interface {
	CheckSignIncomingKey([]byte, []byte, []byte) error
	DecryptIncomingKey([]byte, []byte, []byte) []byte
}

type EncryptionKey interface {
	EncryptAes([]byte, []byte) ([]byte, error)
}

type ConverterKey interface {
	JsonConverter(any) ([]byte, error)
	ConverterToPrivateKey([]byte) (*rsa.PrivateKey, error)
	ConvertDataToHash([]byte, []byte) []byte
}
