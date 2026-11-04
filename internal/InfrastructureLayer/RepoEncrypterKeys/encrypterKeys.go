package RepoEncrypterKeys

import (
	"Kaban/internal/DomainLevel"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"sync"

	"github.com/awnumar/memguard"
)

type (
	keys struct {
		oldKey *memguard.LockedBuffer
		newKey *memguard.LockedBuffer
		rm     *sync.RWMutex
	}
)

func SetKeys() {
	r, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	DomainLevel.GetCryptoKeys = &keys{
		oldKey: &memguard.LockedBuffer{},
		newKey: &memguard.LockedBuffer{},
		rm:     &sync.RWMutex{},
	}
	DomainLevel.GetCryptoKeys.UpdateNewKey(x509.MarshalPKCS1PrivateKey(r))
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
func (s *keys) UpdateNewKey(key []byte) {

	s.rm.Lock()
	if s.newKey != nil {
		s.newKey.Destroy()
	}
	s.newKey = memguard.NewBufferFromBytes(key)
	s.rm.Unlock()
	return
}
func (s *keys) UpdateOldKey() {
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
