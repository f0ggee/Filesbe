package Handlers

import (
	"context"
	"crypto/sha256"
	"log/slog"
	"sync"

	"github.com/awnumar/memguard"
)

var Keys struct {
	Mut           sync.RWMutex
	NewPrivateKey *memguard.LockedBuffer
	OldPrivateKey *memguard.LockedBuffer
}

func ChangerOldKey() {
	Keys.Mut.Lock()
	Keys.OldPrivateKey.Destroy()
	Keys.OldPrivateKey = memguard.NewBuffer(Keys.NewPrivateKey.Size())
	Keys.OldPrivateKey.Copy(Keys.NewPrivateKey.Bytes())
	Keys.Mut.Unlock()
}

func (sa *HandlerPackCollect) SwapKeys() bool {

	slog.Info("SwapKeys", "Start", true)
	defer Keys.NewPrivateKey.Destroy()
	ChangerOldKey()

	aesKey, plaintext, sign, err := sa.RedisControlling.Reader.GetKey(context.Background())
	if err != nil {
		return false
	}

	hashSha := sha256.New()
	hashSha.Write(plaintext)
	hashSha.Write(aesKey)

	err = sa.Crypto.Validate.CheckSignKey(sign, hashSha.Sum([]byte(nil)), ControlPrivateKeyStruct.MasterServerPublicKeyBytes)
	if err != nil {
		slog.Error("Error checkSignIncomingKey", "Error", err.Error())
		return false
	}

	AesKeyDecrypted1, err2 := sa.Crypto.Decrypt.DecryptAesKey(aesKey, ControlPrivateKeyStruct.OurPrivateKeyIntoBytes)
	if err2 != nil {
		return false
	}

	NewRsaKey := (sa.Crypto.Decrypt.DecryptPacket(AesKeyDecrypted1, plaintext))
	if NewRsaKey == nil {
		return false
	}
	defer NewRsaKey.Destroy()

	Keys.Mut.Lock()
	Keys.NewPrivateKey = memguard.NewBuffer(NewRsaKey.Size())
	Keys.NewPrivateKey.Copy(NewRsaKey.Bytes())
	Keys.Mut.Unlock()
	slog.Info("SwapKeys", "End", true)
	return true
}
