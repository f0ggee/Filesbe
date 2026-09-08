package DomainLevel

import (
	"crypto/rsa"
	"io"

	"github.com/awnumar/memguard"
)

type Decryption interface {
	DecryptPacket([]byte, []byte) *memguard.LockedBuffer
	DecryptAesKey([]byte, []byte) ([]byte, error)
	DecryptFileInfo([]byte, []byte, []byte) ([]byte, string, error)
	SayHello(string) string
}

type IncomeData struct {
	Key  []byte
	Data []byte
}
type Decrypter interface {
	DecryptData(IncomeData) ([]byte, error)
}
type FileLabelsBytes struct {
	FileName string
	AesKey   string
}

type CheckSignKeyIncomingData struct {
	Sign            []byte
	Hash            []byte
	MasterPublicKey []byte
}
type CryptoValidating interface {
	CheckSign(CheckSignKeyIncomingData) error
	PasswordVerify([]byte, []byte) error
}

type CryptoGenerating interface {
	GenerateText(int) string
	GenerateSignature(message []byte, key []byte) ([]byte, error)
	GenerateHash([]byte) ([]byte, error)
}
type Encryption interface {
	EncryptAes([]byte, []byte) ([]byte, error)
	EncryptFileInfo([]byte, *rsa.PublicKey) ([]byte, error)
}

type Crypto interface {
	Encrypter([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}
type StreamCrypto interface {
	EncryptData(io.Reader) ([]byte, error)
	DecryptData(io.Reader) error
}

type MakeCrypto interface {
	InitializerCrypto([]byte) Crypto
}
type MakeStreamCrypto interface {
	InitializerStreamCrypto([]byte) StreamCrypto
}
