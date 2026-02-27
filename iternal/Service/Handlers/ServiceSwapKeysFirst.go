package Handlers

import (
	"log/slog"

	"github.com/awnumar/memguard"
)

func (sa *HandlerPackCollect) SwapKeyFirst() {

	slog.Info("SwapKeyFirst", "Start", true)

	aesKey, plainText, sign, err := sa.S.RedisConn.GetKey()
	if err != nil {
		return
	}

	defer memguard.WipeBytes(aesKey)
	defer memguard.WipeBytes(plainText)

	shaHashFromData := sa.S.Choose.ConvertDataToHash(plainText, aesKey)

	err = sa.S.Choose.CheckSignIncomingKey(sign, shaHashFromData, ControlPrivateKeyStruct.MasterServerPublicKeyBytes)
	if err != nil {
		slog.Error("Error check sign incomingKey", "Error", err.Error())
		return
	}

	NewRsaKey := memguard.NewBufferFromBytes(sa.S.Choose.DecryptIncomingKey(aesKey, plainText, ControlPrivateKeyStruct.OurPrivateKeyIntoBytes))
	if NewRsaKey == nil {
		return
	}
	defer NewRsaKey.Destroy()

	Keys.NewPrivateKey = memguard.NewBuffer(NewRsaKey.Size())
	Keys.OldPrivateKey = memguard.NewBuffer(NewRsaKey.Size())
	Keys.NewPrivateKey.Copy(NewRsaKey.Bytes())
	Keys.OldPrivateKey.Copy(NewRsaKey.Bytes())
	slog.Info("SwapKeyFirst", "End", true)

}
