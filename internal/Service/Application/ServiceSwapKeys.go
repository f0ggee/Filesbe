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
	GetCrypto
	GetControlKeys
	EncrypterKeys
	Parser
}

func GetNewNewSwapKeys(redisControlling RedisControlling, crypto GetCrypto, controlKeys GetControlKeys, encrypterKeys EncrypterKeys, parser Parser) *NewSwapKeys {
	return &NewSwapKeys{RedisControlling: redisControlling, GetCrypto: crypto, GetControlKeys: controlKeys, EncrypterKeys: encrypterKeys, Parser: parser}
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

	gerOurPrivateKey, err := sa.GetControlKeys.Keys.GerOurPrivateKey()
	if err != nil {
		return 0
	}
	AesKeyDecrypted1, err2 := sa.GetCrypto.Decrypt.DecryptAesKey(gerOurPrivateKey, grpcPacket.AesKey)
	if err2 != nil {
		return DomainLevel.DefaultErrorTime
	}
	NewRsaKey := sa.GetCrypto.Decrypt.DecryptPacket(AesKeyDecrypted1, grpcPacket.PlainText)
	if NewRsaKey == nil {
		return DomainLevel.DefaultErrorTime
	}
	defer NewRsaKey.Destroy()

	hashSha := sha256.New()
	hashSha.Write(NewRsaKey.Bytes())

	getMasterPublicKey, err := sa.GetControlKeys.Keys.GetMasterPublicKey()
	if err != nil {
		return DomainLevel.DefaultErrorTime
	}

	err = sa.GetCrypto.Validate.CheckSignKey(DomainLevel.CheckSignKeyIncomingData{
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
