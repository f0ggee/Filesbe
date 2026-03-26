package KeysManager

import (
	"crypto/rand"
	"log/slog"
	"sync"

	"github.com/awnumar/memguard"
)

type Updater struct {
	Mu            *sync.RWMutex
	NewPrivateKey *memguard.LockedBuffer
	OldPrivateKey *memguard.LockedBuffer
}

func (u Updater) UpdateKey(key *memguard.LockedBuffer) {
	defer u.Mu.Unlock()
	u.Mu.Lock()
	u.NewPrivateKey = memguard.NewBuffer(key.Size())
	u.NewPrivateKey.Copy(key.Bytes())
	slog.Info("Updating key", u.NewPrivateKey.String())

	return
}

func (u Updater) GetKey() []byte {
	defer u.Mu.RUnlock()

	u.Mu.RLock()

	slog.Info("GetKey", u.NewPrivateKey.String())
	return u.NewPrivateKey.Bytes()
}

func (u Updater) GetOldKey() []byte {

	slog.Info("GetOldKey", u.NewPrivateKey.String())
	return u.OldPrivateKey.Bytes()
}

func (u Updater) UpdateOldKey() {
	defer u.Mu.Unlock()
	u.Mu.Lock()

	u.OldPrivateKey.Destroy()
	u.OldPrivateKey = memguard.NewBuffer(u.NewPrivateKey.Size())
	u.OldPrivateKey.Copy(u.NewPrivateKey.Bytes())
}

func (u Updater) FillOldKey() {

	defer u.Mu.Unlock()
	u.Mu.Lock()
	err := error(nil)
	u.NewPrivateKey, err = memguard.NewBufferFromReader(rand.Reader, 32)
	if err != nil {
		slog.Error("Error filling old key", err.Error())
		return
	}
	u.OldPrivateKey, err = memguard.NewBufferFromReader(rand.Reader, 32)
	if err != nil {
		slog.Error("Error filling old key", err.Error())
		return
	}

}
