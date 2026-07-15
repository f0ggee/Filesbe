package DomainLevel

const (
	ErrorKeyNotCreated = "the key wasn't created"
)

type NewServerKeys struct {
	masterKeyPublic  []byte
	workerPrivateKey []byte
}

func (s NewServerKeys) GerOurPrivateKey() []byte {
	return s.workerPrivateKey
}

func GetNewSetKeys(ourPrivateKeyIntoBytes []byte, masterServerPublicKeyBytes []byte) *NewServerKeys {
	return &NewServerKeys{workerPrivateKey: ourPrivateKeyIntoBytes, masterKeyPublic: masterServerPublicKeyBytes}
}

func (s NewServerKeys) GetMasterPublicKey() []byte {
	return s.masterKeyPublic
}
