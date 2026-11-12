package RepoEncrypterKeys

import (
	"Kaban/internal/DomainLevel"
	"crypto/rand"
	"encoding/hex"
	"os"
	"sync"

	"github.com/awnumar/memguard"
	"github.com/joho/godotenv"
)

type (
	keys struct {
		oldKey    *memguard.LockedBuffer
		newKey    *memguard.LockedBuffer
		rm        *sync.RWMutex
		masterKey []byte
		serverKey []byte
	}
)

func (s *keys) GetServerKey() []byte {

	return s.serverKey
}

func (s *keys) GetMasterKey() []byte {
	return s.masterKey
}

func SetKeys() {

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
	DomainLevel.GetCryptoKeys = &keys{
		oldKey:    &memguard.LockedBuffer{},
		newKey:    &memguard.LockedBuffer{},
		rm:        &sync.RWMutex{},
		masterKey: k1,
		serverKey: k2,
	}
	DomainLevel.GetCryptoKeys.UpdateNewKey([]byte(rand.Text()))
}
func (s *keys) GetKey() []byte {
	s.rm.RLock()
	defer s.rm.RUnlock()
	return s.newKey.Data()
}
func (s *keys) GetOldKey() []byte {
	s.rm.RLock()
	defer s.rm.RUnlock()
	return s.oldKey.Bytes()
}
func (s *keys) UpdateNewKey(newKey []byte) {

	s.rm.Lock()
	if s.newKey != nil {
		s.newKey.Destroy()
	}
	s.newKey = memguard.NewBufferFromBytes(newKey)
	s.rm.Unlock()
	return
}
func (s *keys) UpdateOldKey([]byte) {
	s.rm.Lock()
	if s.oldKey != nil {
		s.oldKey.Destroy()
	}

	if s.newKey != nil {
		s.oldKey = memguard.NewBuffer(s.newKey.Size())
		s.oldKey.Copy(s.newKey.Bytes())
	}

	s.rm.Unlock()
}
