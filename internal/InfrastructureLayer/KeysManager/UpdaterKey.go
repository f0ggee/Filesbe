package KeysManager

import (
	"github.com/awnumar/memguard"
)

func (u *Updater) UpdateKey(key *memguard.LockedBuffer) {
	u.Mu.Lock()

	u.NewPrivateKey.Destroy()
	u.NewPrivateKey = memguard.NewBuffer(key.Size())
	u.NewPrivateKey.Copy(key.Data())

	u.Mu.Unlock()
	return
}

func (u Updater) GetKey2() []byte {

	return u.NewPrivateKey.Bytes()
}
