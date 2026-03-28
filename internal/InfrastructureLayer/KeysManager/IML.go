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
