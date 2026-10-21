package DomainLevel

type IncomeData struct {
	Key  []byte
	Data []byte
}

type FileLabelsBytes struct {
	FileName string
	AesKey   []byte
}

type CheckSignKeyIncomingData struct {
	Sign            []byte
	Hash            []byte
	MasterPublicKey []byte
}
type CryptoValidating interface {
	PasswordVerify([]byte, []byte) error
}

type CryptoGenerating interface {
	GenerateText(int) string
	GenerateHash([]byte) ([]byte, error)
}
type Signer interface {
	MakeSign(signData []byte) ([]byte, error)
	CheckSign(signature []byte) error
}
type SignatureMaker interface {
	SetAdditionalData(data []byte) SignatureMaker
	MakeSigner(key []byte) (Signer, error)
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
