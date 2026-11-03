package RepoEncrypterKeys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"sync"

	"github.com/awnumar/memguard"
)

type Keys struct {
	oldKey *memguard.LockedBuffer
	newKey *memguard.LockedBuffer
	rm     *sync.RWMutex
}

var KeysControl Keys

func init() {
	keys, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	KeysControl = Keys{
		oldKey: &memguard.LockedBuffer{},
		newKey: &memguard.LockedBuffer{},
		rm:     &sync.RWMutex{},
	}
	x := x509.MarshalPKCS1PrivateKey(keys)
	KeysControl.UpdateNewKey(x)

}

func GetKeys() Keys {

	var k Keys
	k.newKey, _ = memguard.NewBufferFromReader(rand.Reader, 32)
	k.oldKey, _ = memguard.NewBufferFromReader(rand.Reader, 32)
	return k
}

func (s *Keys) GetKey() []byte {
	s.rm.RLock()
	defer s.rm.RUnlock()
	return s.newKey.Data()
}
func (s *Keys) GetOldKey() []byte {
	s.rm.RLock()
	defer s.rm.RUnlock()
	return s.oldKey.Bytes()
}
func (s *Keys) UpdateNewKey(key []byte) {

	s.rm.Lock()
	if s.newKey != nil {
		s.newKey.Destroy()
	}
	s.newKey = memguard.NewBuffer(len(key))
	s.newKey.Copy(key)
	s.rm.Unlock()
	return
}
func (s *Keys) UpdateOldKey() {
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
