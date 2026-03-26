package Handlers

import (
	"context"
	"crypto/sha256"
	"log/slog"
)

func (sa *HandlerPackCollect) SwapKeys() bool {

	slog.Info("SwapKeys", "Start", true)
	sa.Keys.ControllerKey.UpdateOldKey()
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

	NewRsaKey := sa.Crypto.Decrypt.DecryptPacket(AesKeyDecrypted1, plaintext)
	if NewRsaKey == nil {
		return false
	}
	defer NewRsaKey.Destroy()

	sa.Keys.ControllerKey.UpdateKey(NewRsaKey)

	slog.Info("SwapKeys", "End", true)
	return true
}
