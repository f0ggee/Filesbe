package DomainLevel

type IncomeData struct {
	Key  []byte
	Data []byte
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

type Crypto interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

type CryptoMaker interface {
	// MakeCrypto ACD- additional crypto data
	MakeCrypto(key []byte, ACD []byte) (Crypto, error)
	GetRequiredOverheadSize() int
}
