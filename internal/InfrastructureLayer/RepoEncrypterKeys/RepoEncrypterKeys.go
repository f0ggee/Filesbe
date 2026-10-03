package RepoEncrypterKeys

import (
	"github.com/awnumar/memguard"
)

const (
	ErrorUpdateNewKey = "OldKey needs to be updated first"
)

type Keys struct {
	OldKey         *memguard.LockedBuffer
	NewKey         *memguard.LockedBuffer
	isOldKeyCreate bool
}

func GetNewKeys() Keys {
	return Keys{OldKey: &memguard.LockedBuffer{}, NewKey: &memguard.LockedBuffer{}}
}
func (s *Keys) GetKey() []byte    { return s.NewKey.Data() }
func (s *Keys) GetOldKey() []byte { return s.OldKey.Bytes() }
func (s *Keys) UpdateNewKey(key []byte) {
	if s.NewKey != nil {
		s.NewKey.Destroy()
	}
	s.NewKey = memguard.NewBuffer(len(key))
	s.NewKey.Copy(key)
	return
}
func (s *Keys) UpdateOldKey() {
	if s.OldKey != nil {
		s.OldKey.Destroy()
	}

	if s.NewKey != nil {
		s.OldKey = memguard.NewBuffer(s.NewKey.Size())
		s.OldKey.Copy(s.NewKey.Bytes())
	}

}
