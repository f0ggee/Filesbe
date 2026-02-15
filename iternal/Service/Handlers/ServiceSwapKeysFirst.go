package Handlers

import (
	"Kaban/iternal/InfrastructureLayer"
	"log/slog"

	"github.com/awnumar/memguard"
)

func SwapKeyFirst() {

	slog.Info("SwapKeyFirst", "Start", true)
	redis := *InfrastructureLayer.NewSetRedisConnect()
	encryptionKeyS := *InfrastructureLayer.ConnectToEncryptKey()

	aesKey, plainText, sign, err := redis.Ras.GetKey()
	if err != nil {
		return
	}

	defer memguard.WipeBytes(aesKey)
	defer memguard.WipeBytes(plainText)

	shaHashFromData := encryptionKeyS.Choose.ConvertDataToHash(plainText, aesKey)

	err = encryptionKeyS.Choose.CheckSignIncomingKey(sign, shaHashFromData, ControlPrivateKeyStruct.MasterServerPublicKeyBytes)
	if err != nil {
		slog.Error("Error check sign incomingKey", "Error", err.Error())
		return
	}

	NewRsaKey := memguard.NewBufferFromBytes(encryptionKeyS.Choose.DecryptIncomingKey(aesKey, plainText, ControlPrivateKeyStruct.OurPrivateKeyIntoBytes))
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
