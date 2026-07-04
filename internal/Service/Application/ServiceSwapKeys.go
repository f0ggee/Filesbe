package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"context"
	"crypto/sha256"
	"time"
)

type NewSwapKeys struct {
	RedisControlling
	Crypto
	ControlKeys
	EncrypterKeys
	Parser
}

func GetNewNewSwapKeys(redisControlling RedisControlling, crypto Crypto, controlKeys ControlKeys, encrypterKeys EncrypterKeys, parser Parser) *NewSwapKeys {
	return &NewSwapKeys{RedisControlling: redisControlling, Crypto: crypto, ControlKeys: controlKeys, EncrypterKeys: encrypterKeys, Parser: parser}
}

func (sa *NewSwapKeys) SetSwapKeys() time.Duration {
	Data, err := sa.RedisControlling.Reader.GetKey(context.Background())
	if err != nil {
		return DomainLevel.DefaultErrorTime
	}
	grpcPacket := &Dto.RedisPacketStructFromMasterServer{
		AesKey:          nil,
		PlainText:       nil,
		Signature:       nil,
		TimeNextSwaping: time.Duration(0),
	}

	err = sa.Parser.Decode.JsonDecodeMarshall(&grpcPacket, Data)
	if err != nil {
		return DomainLevel.DefaultErrorTime
	}

	gerOurPrivateKey, err := sa.ControlKeys.Keys.GerOurPrivateKey()
	if err != nil {
		return 0
	}
	AesKeyDecrypted1, err2 := sa.Crypto.Decrypt.DecryptAesKey(gerOurPrivateKey, grpcPacket.AesKey)
	if err2 != nil {
		return DomainLevel.DefaultErrorTime
	}
	NewRsaKey := sa.Crypto.Decrypt.DecryptPacket(AesKeyDecrypted1, grpcPacket.PlainText)
	if NewRsaKey == nil {
		return DomainLevel.DefaultErrorTime
	}
	defer NewRsaKey.Destroy()

	hashSha := sha256.New()
	hashSha.Write(NewRsaKey.Bytes())

	getMasterPublicKey, err := sa.ControlKeys.Keys.GetMasterPublicKey()
	if err != nil {
		return DomainLevel.DefaultErrorTime
	}

	err = sa.Crypto.Validate.CheckSignKey(DomainLevel.CheckSignKeyIncomingData{
		Sign:            grpcPacket.Signature,
		Hash:            hashSha.Sum(nil),
		MasterPublicKey: getMasterPublicKey,
	})
	if err != nil {
		return DomainLevel.DefaultErrorTime
	}

	sa.EncrypterKeys.GetKeys.UpdateOldKey()
	err = sa.EncrypterKeys.GetKeys.UpdateNewKey(NewRsaKey)
	if err != nil {
		return DomainLevel.DefaultErrorTime
	}
	return grpcPacket.TimeNextSwaping
}
