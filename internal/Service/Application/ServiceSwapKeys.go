package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"context"
	"crypto/sha256"
	"time"
)

type NewSwapKeys struct {
	redisControlling
	getCrypto
	getControlKeys
	encrypterKeys
	parser
}

func GetNewNewSwapKeys(redisControlling redisControlling, crypto getCrypto, controlKeys getControlKeys, encrypterKeys encrypterKeys, parser parser) *NewSwapKeys {
	return &NewSwapKeys{redisControlling: redisControlling, getCrypto: crypto, getControlKeys: controlKeys, encrypterKeys: encrypterKeys, parser: parser}
}

func (sa *NewSwapKeys) SetSwapKeys() time.Duration {
	Data, err := sa.redisControlling.Reader.GetKey(context.Background())
	if err != nil {
		return DomainLevel.DefaultErrorTime
	}
	grpcPacket := &Dto.RedisPacketStructFromMasterServer{
		AesKey:          nil,
		PlainText:       nil,
		Signature:       nil,
		TimeNextSwaping: time.Duration(0),
	}

	err = sa.parser.Decode.JsonDecodeMarshall(&grpcPacket, Data)
	if err != nil {
		return DomainLevel.DefaultErrorTime
	}

	AesKeyDecrypted1, err2 := sa.getCrypto.Decrypt.DecryptAesKey(sa.getControlKeys.Keys.GerOurPrivateKey(), grpcPacket.AesKey)
	if err2 != nil {
		return DomainLevel.DefaultErrorTime
	}
	NewRsaKey := sa.getCrypto.Decrypt.DecryptPacket(AesKeyDecrypted1, grpcPacket.PlainText)
	if NewRsaKey == nil {
		return DomainLevel.DefaultErrorTime
	}
	defer NewRsaKey.Destroy()

	hashSha := sha256.New()
	hashSha.Write(NewRsaKey.Bytes())

	getMasterPublicKey := sa.getControlKeys.Keys.GetMasterPublicKey()

	err = sa.getCrypto.Validate.CheckSignKey(DomainLevel.CheckSignKeyIncomingData{
		Sign:            grpcPacket.Signature,
		Hash:            hashSha.Sum(nil),
		MasterPublicKey: getMasterPublicKey,
	})
	if err != nil {
		return DomainLevel.DefaultErrorTime
	}

	sa.encrypterKeys.GetKeys.UpdateOldKey()
	err = sa.encrypterKeys.GetKeys.UpdateNewKey(NewRsaKey)
	if err != nil {
		return DomainLevel.DefaultErrorTime
	}
	return grpcPacket.TimeNextSwaping
}
