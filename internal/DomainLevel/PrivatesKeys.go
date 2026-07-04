package DomainLevel

import "errors"

const (
	ErrorKeyNotCreated = "the key wasn't created"
)

type NewSetKeys struct {
	MasterKeyPublic  []byte
	WorkerPrivateKey []byte
}

func (s NewSetKeys) GerOurPrivateKey() ([]byte, error) {
	if s.WorkerPrivateKey == nil {
		return nil, errors.New(ErrorKeyNotCreated)
	}
	return s.WorkerPrivateKey, nil
}

func GetNewSetKeys(ourPrivateKeyIntoBytes []byte, masterServerPublicKeyBytes []byte) *NewSetKeys {
	return &NewSetKeys{WorkerPrivateKey: ourPrivateKeyIntoBytes, MasterKeyPublic: masterServerPublicKeyBytes}
}

func (s NewSetKeys) GetMasterPublicKey() ([]byte, error) {
	if s.MasterKeyPublic == nil {
		return nil, errors.New(ErrorKeyNotCreated)
	}
	return s.MasterKeyPublic, nil
}
