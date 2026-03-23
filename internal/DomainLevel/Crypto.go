package DomainLevel

import (
	"crypto/rsa"

	"github.com/awnumar/memguard"
)

type Decryption interface {
	DecryptPacket([]byte, []byte) *memguard.LockedBuffer
	DecryptAesKey([]byte, []byte) ([]byte, error)
	DecryptFileInfo([]byte, []byte, []byte) ([]byte, string, error)
}

type CryptoValidating interface {
	CheckSignKey([]byte, []byte, []byte) error
	//CheckSignatureGrpc([]byte, []byte, []byte) error
}
type CryptoGenerating interface {
	GenerateShortName() string
	GenerateSignature(message []byte, key []byte) ([]byte, error)
}
type Encryption interface {
	EncryptAes([]byte, []byte) ([]byte, error)
	EncryptFileInfo([]byte, *rsa.PublicKey) ([]byte, error)
}
