package RepoEncrypterKeys

import (
	"sync"

	"github.com/awnumar/memguard"
)

type Keys struct {
	OldKey *memguard.LockedBuffer
	NewKey *memguard.LockedBuffer
	Rm     *sync.RWMutex
}

func GetKeys() Keys {
	return Keys{
		OldKey: &memguard.LockedBuffer{},
		NewKey: &memguard.LockedBuffer{},
		Rm:     &sync.RWMutex{},
	}
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
