package DomainLevel

const (
	ErrorKeyNotCreated = "the key wasn't created"
)

type NewSetKeys struct {
	masterKeyPublic  []byte
	workerPrivateKey []byte
}

func (s NewSetKeys) GerOurPrivateKey() []byte {
	return s.workerPrivateKey
}

func GetNewSetKeys(ourPrivateKeyIntoBytes []byte, masterServerPublicKeyBytes []byte) *NewSetKeys {
	return &NewSetKeys{workerPrivateKey: ourPrivateKeyIntoBytes, masterKeyPublic: masterServerPublicKeyBytes}
}

func (s NewSetKeys) GetMasterPublicKey() []byte {
	return s.masterKeyPublic
}
