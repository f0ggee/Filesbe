package RepoEncrypterKeys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"sync"

	"github.com/awnumar/memguard"
)

type Keys struct {
	OldKey *memguard.LockedBuffer
	NewKey *memguard.LockedBuffer
	Rm     *sync.RWMutex
}

func GetKeysTest() Keys {
	keys, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	var k Keys
	x := x509.MarshalPKCS1PrivateKey(keys)
	k.Rm = &sync.RWMutex{}
	k.UpdateNewKey(x)
	return k
}
func GetKeys() Keys {

	var k Keys
	k.NewKey, _ = memguard.NewBufferFromReader(rand.Reader, 32)
	k.OldKey, _ = memguard.NewBufferFromReader(rand.Reader, 32)
	return k
}

func (s *Keys) GetKey() []byte {
	return s.NewKey.Data()
}
func (s *Keys) GetOldKey() []byte { return s.OldKey.Bytes() }
func (s *Keys) UpdateNewKey(key []byte) {

	s.Rm.Lock()
	if s.NewKey != nil {
		s.NewKey.Destroy()
	}
	s.NewKey = memguard.NewBuffer(len(key))
	s.NewKey.Copy(key)
	s.Rm.Unlock()
	return
}
func (s *Keys) UpdateOldKey() {
	s.Rm.Lock()
	if s.OldKey != nil {
		s.OldKey.Destroy()
	}

	if s.NewKey != nil {
		s.OldKey = memguard.NewBuffer(s.NewKey.Size())
		s.OldKey.Copy(s.NewKey.Bytes())
	}

	s.Rm.Unlock()
}
