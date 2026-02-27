package Handlers

import (
	"log/slog"
	"sync"

	"github.com/awnumar/memguard"
)

var Keys struct {
	Mut           sync.RWMutex
	NewPrivateKey *memguard.LockedBuffer
	OldPrivateKey *memguard.LockedBuffer
}

//SwapKeys generates a  pair keys

func (sa *HandlerPackCollect) SwapKeys() bool {

	slog.Info("SwapKeys", "Start", true)
	Keys.Mut.Lock()

	Keys.OldPrivateKey.Destroy()
	Keys.OldPrivateKey = memguard.NewBuffer(Keys.NewPrivateKey.Size())
	Keys.OldPrivateKey.Copy(Keys.NewPrivateKey.Bytes())
	Keys.NewPrivateKey.Destroy()

	aesKey, plaintext, sign, err := sa.S.RedisConn.GetKey()
	if err != nil {
		return false
	}

	hashFromData := memguard.NewBufferFromBytes(sa.S.Choose.ConvertDataToHash(plaintext, aesKey))
	if hashFromData == nil {
		return false
	}
	defer hashFromData.Destroy()

	err = sa.S.Choose.CheckSignIncomingKey(sign, hashFromData.Bytes(), ControlPrivateKeyStruct.MasterServerPublicKeyBytes)
	if err != nil {
		slog.Error("Error checkSignIncomingKey", "Error", err.Error())
		return false
	}
	NewRsaKey := memguard.NewBufferFromBytes(sa.S.Choose.DecryptIncomingKey(aesKey, plaintext, ControlPrivateKeyStruct.OurPrivateKeyIntoBytes))
	if NewRsaKey == nil {
		return false
	}
	defer NewRsaKey.Destroy()
	Keys.NewPrivateKey = memguard.NewBuffer(NewRsaKey.Size())
	Keys.NewPrivateKey.Copy(NewRsaKey.Bytes())
	slog.Info("SwapKeys", "End", true)

	Keys.Mut.Unlock()
	return true
}
