package DomainLevel

import "os"

type NewServerKeys struct {
	masterKeyPublic  []byte
	workerPrivateKey []byte
}

var ServerKeys NewServerKeys

func (s NewServerKeys) GerOurPrivateKey() []byte {
	return s.workerPrivateKey
}

func init() {
	masterKey := os.Getenv("Public_Key_Master_Server")
	ourKey := os.Getenv("Our_Private_Key")
	ServerKeys = NewServerKeys{
		masterKeyPublic:  []byte(masterKey),
		workerPrivateKey: []byte(ourKey),
	}
}

func (s NewServerKeys) GetMasterPublicKey() []byte {
	return s.masterKeyPublic
}
