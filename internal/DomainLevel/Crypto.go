package DomainLevel

type IncomeData struct {
	Key  []byte
	Data []byte
}

type FileLabelsBytes struct {
	FileName []byte
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

type (
	CryptoMaker interface {
		// MakeCrypto ACD- additional crypto data
		MakeCrypto(key []byte, ACD []byte) (Crypto, error)
		GetRequiredOverheadSize() int
	}
	Crypto interface {
		Encrypt([]byte) ([]byte, error)
		Decrypt([]byte) ([]byte, error)
	}
)
type (
	CryptoKeys interface {
		GetKey() []byte
		GetOldKey() []byte
		UpdateNewKey(newKey []byte)
		//UpdateOldKey uses a new key as an old key by default but if we have to pass another key then we can directly
		//use it in a method parameter. However, if we want to do planning exchange, we don't need to pass anything.
		UpdateOldKey(otherKey []byte)
	}
)

var GetCryptoKeys CryptoKeys
