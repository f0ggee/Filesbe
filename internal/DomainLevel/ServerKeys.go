package DomainLevel

import (
	"os"

	"github.com/joho/godotenv"
)

type NewServerKeys struct {
	MasterKeyPublic  []byte
	WorkerPrivateKey []byte
}

var ServerKeys NewServerKeys

func (s NewServerKeys) GerOurPrivateKey() []byte {
	return s.WorkerPrivateKey
}

func init() {
	err := godotenv.Load()
	if err != nil {
		//panic(err)
	}
	masterKey := os.Getenv("Public_Key_Master_Server")
	ourKey := os.Getenv("Our_Private_Key")
	ServerKeys = NewServerKeys{
		MasterKeyPublic:  []byte(masterKey),
		WorkerPrivateKey: []byte(ourKey),
	}
}

func (s NewServerKeys) GetMasterPublicKey() []byte {
	return s.MasterKeyPublic
}
