package Handlers

import (
	"Kaban/internal/Dto"
	"context"
	"crypto/sha256"
	"encoding/json"
	"log/slog"
)

func ConvertData(Data []byte) (*Dto.RedisPacketStructFromMasterServer, error) {
	MasterDataLooks := &Dto.RedisPacketStructFromMasterServer{
		AesKey:    nil,
		PlainText: nil,
		Signature: nil,
	}
	err := json.Unmarshal(Data, &MasterDataLooks)
	if err != nil {
		slog.Error("SwapKeys", "Unmarshal err", err)
		return nil, err
	}
	return MasterDataLooks, nil
}

func (sa *HandlerPackCollect) SwapKeys() bool {
	slog.Info("SwapKeys", "Start", true)

	Data, err := sa.RedisControlling.Reader.GetKey(context.Background())
	if err != nil {
		return false
	}

	UnpackedData, err := ConvertData(Data)
	if err != nil {
		return false
	}

	AesKeyDecrypted1, err2 := sa.Crypto.Decrypt.DecryptAesKey(sa.Keys.ControllerKey.GetOurKey(), UnpackedData.AesKey)
	if err2 != nil {
		return false
	}

	NewRsaKey := sa.Crypto.Decrypt.DecryptPacket(AesKeyDecrypted1, UnpackedData.PlainText)
	if NewRsaKey == nil {
		return false
	}
	defer NewRsaKey.Destroy()

	hashSha := sha256.New()
	hashSha.Write(NewRsaKey.Bytes())

	err = sa.Crypto.Validate.CheckSignKey(UnpackedData.Signature, hashSha.Sum([]byte(nil)), sa.Keys.ControllerKey.GetMasterKey())
	if err != nil {
		return false
	}

	sa.Keys.ControllerKey.UpdateOldKey()
	sa.Keys.ControllerKey.UpdateKey(NewRsaKey)
	slog.Info("SwapKeys", "End", true)
	return true
}
