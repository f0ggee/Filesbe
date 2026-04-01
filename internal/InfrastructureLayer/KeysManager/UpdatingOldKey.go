package KeysManager

import "github.com/awnumar/memguard"

func (u *Updater) UpdateOldKey() {
	defer u.Mu.Unlock()
	u.Mu.Lock()

	u.OldPrivateKey.Destroy()
	u.OldPrivateKey = memguard.NewBuffer(u.NewPrivateKey.Size())
	u.OldPrivateKey.Copy(u.NewPrivateKey.Bytes())
}
