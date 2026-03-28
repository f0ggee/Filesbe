package KeysManager

import "github.com/awnumar/memguard"

func (u Updater) UpdateKey(key *memguard.LockedBuffer) {
	defer u.Mu.Unlock()
	u.Mu.Lock()
	u.NewPrivateKey = memguard.NewBuffer(key.Size())
	u.NewPrivateKey.Copy(key.Bytes())

	return
}
