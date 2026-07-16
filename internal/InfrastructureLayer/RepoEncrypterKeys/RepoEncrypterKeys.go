package RepoEncrypterKeys

import (
	"errors"

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

func GetNewKeys(oldKey *memguard.LockedBuffer, newKey *memguard.LockedBuffer) *Keys {

	oldKey = memguard.NewBufferRandom(32)
	newKey = memguard.NewBufferRandom(32)
	return &Keys{OldKey: oldKey, NewKey: newKey}
}
func (s *Keys) GetKey() []byte    { return s.NewKey.Data() }
func (s *Keys) GetOldKey() []byte { return s.OldKey.Bytes() }
func (s *Keys) UpdateNewKey(key *memguard.LockedBuffer) error {
	if !s.isOldKeyCreate {
		return errors.New(ErrorUpdateNewKey)
	}
	s.isOldKeyCreate = false
	s.NewKey.Destroy()
	s.NewKey = memguard.NewBuffer(key.Size())
	s.NewKey.Copy(key.Data())
	return nil
}
func (s *Keys) UpdateOldKey() {
	defer func() {
		s.isOldKeyCreate = true
	}()
	s.OldKey.Destroy()
	s.OldKey = memguard.NewBuffer(s.NewKey.Size())
	s.OldKey.Copy(s.NewKey.Bytes())
}
