package DomainLevel

import (
	"crypto/rsa"

	"github.com/awnumar/memguard"
)

type KeyInteraction interface {
	CheckSignIncomingKey([]byte, []byte, []byte) error
	DecryptPacket([]byte, []byte) *memguard.LockedBuffer
	DecryptAesKey([]byte, []byte) ([]byte, error)
}

type EncryptionKey interface {
	EncryptAes([]byte, []byte) ([]byte, error)
}

type Converter interface {
	JsonConverter(any) ([]byte, error)
	ConverterToPrivateKey([]byte) (*rsa.PrivateKey, error)
	ConverterToPublicKey([]byte) (*rsa.PublicKey, error)
	ConvertDataToHash([]byte, []byte) []byte
}
