package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"context"
	"crypto/sha256"
	"encoding/json"
	"log/slog"
	"time"
)

func ConvertData(Data []byte) (*Dto.RedisPacketStructFromMasterServer, error) {
	MasterDataLooks := &Dto.RedisPacketStructFromMasterServer{
		AesKey:          nil,
		PlainText:       nil,
		Signature:       nil,
		TimeNextSwaping: time.Duration(0),
	}
	err := json.Unmarshal(Data, &MasterDataLooks)
	if err != nil {
		slog.Error("Func ConvertData: Error", "Unmarshal err", err)
		return nil, err
	}
	return MasterDataLooks, nil
}

func (sa *HandlerPackCollect) SwapKeys() time.Duration {
	slog.Info("Func SwapKeys:", "Start", true)
	Data, err := sa.RedisControlling.Reader.GetKey(context.Background())
	if err != nil {
		return DomainLevel.DefaultErrorTime
	}

	UnpackedData, err := ConvertData(Data)
	if err != nil {
		return DomainLevel.DefaultErrorTime
	}

	AesKeyDecrypted1, err2 := sa.Crypto.Decrypt.DecryptAesKey(sa.Keys.ControllerKey.GetOurKey(), UnpackedData.AesKey)
	if err2 != nil {
		return DomainLevel.DefaultErrorTime
	}
	NewRsaKey := sa.Crypto.Decrypt.DecryptPacket(AesKeyDecrypted1, UnpackedData.PlainText)
	if NewRsaKey == nil {
		return DomainLevel.DefaultErrorTime
	}
	defer NewRsaKey.Destroy()

	hashSha := sha256.New()
	hashSha.Write(NewRsaKey.Bytes())

	err = sa.Crypto.Validate.CheckSignKey(UnpackedData.Signature, hashSha.Sum([]byte(nil)), sa.Keys.ControllerKey.GetMasterKey())
	if err != nil {
		return DomainLevel.DefaultErrorTime
	}

	sa.Keys.ControllerKey.UpdateOldKey()
	sa.Keys.ControllerKey.UpdateKey(NewRsaKey)
	slog.Info("Func SwapKeys:", slog.Group("Data about the exchange", slog.Duration("Time for next swaping", UnpackedData.TimeNextSwaping)))
	return UnpackedData.TimeNextSwaping
}
