package DomainLevel

import (
	"crypto/rsa"

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
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

type CryptoMaker interface {
	// MakeCrypto ACD- additional crypto data
	MakeCrypto(key []byte, ACD []byte) (Crypto, error)
	GetRequiredRandomSize() int
}
