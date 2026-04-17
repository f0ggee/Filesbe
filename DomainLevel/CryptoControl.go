package DomainLevel

type Encryption interface {
	EncryptRsaKey([]byte, []byte) ([]byte, error)
	EncryptAesKey([]byte, []byte) ([]byte, error)
}

type CryptoGenerator interface {
	SignerData([]byte, []byte) ([]byte, error)
	GenerateHash([]byte, []byte) []byte
	GrpcSignerKey() ([]byte, error)
}

type CryptoKeyManager interface {
	GetMasterKey() []byte
}
type Decryptor interface {
	DecrypterCipherData([]byte, []byte) ([]byte, error)
	GrpcDecrypterAesKey([]byte) ([]byte, error)
	SayHi() string
}
