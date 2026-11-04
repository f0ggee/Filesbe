package DomainLevel

import (
	"encoding/hex"
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
		panic(err)
	}
	masterKey := os.Getenv("Public_Key_Master_Server")
	ourKey := os.Getenv("Our_Private_Key")
	k1, err := hex.DecodeString(masterKey)
	if err != nil {
		panic(err)
	}
	k2, err := hex.DecodeString(ourKey)
	if err != nil {
		panic(err)
	}
	ServerKeys = NewServerKeys{
		MasterKeyPublic:  k1,
		WorkerPrivateKey: k2,
	}
}

func (s NewServerKeys) GetMasterPublicKey() []byte {
	return s.MasterKeyPublic
}
