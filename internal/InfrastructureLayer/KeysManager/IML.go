package KeysManager

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"sync"

	"github.com/awnumar/memguard"
)

type Updater struct {
	Mu            *sync.RWMutex
	NewPrivateKey *memguard.LockedBuffer
	OldPrivateKey *memguard.LockedBuffer
	OurPrivateKey string
	MasterKey     string
}

func (u *Updater) GetOurKey() []byte {
	e, err := hex.DecodeString(u.OurPrivateKey)
	if err != nil {

		return nil
	}
	return e
}

func (u *Updater) GetMasterKey() []byte {

	e, err := hex.DecodeString(u.MasterKey)

	if err != nil {
		return nil
	}
	return e
}

func (u *Updater) FillOldKey() {

	defer u.Mu.Unlock()
	u.Mu.Lock()

	err := error(nil)
	u.NewPrivateKey, err = memguard.NewBufferFromReader(rand.Reader, 2)
	if err != nil {
		slog.Error("Error filling new key", err.Error())
		return
	}

	u.OldPrivateKey, err = memguard.NewBufferFromReader(rand.Reader, 32)
	if err != nil {
		slog.Error("Error filling old key", err.Error())
		return
	}

}
