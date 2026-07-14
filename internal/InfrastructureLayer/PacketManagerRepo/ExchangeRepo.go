package PacketManagerRepo

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"context"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/awnumar/memguard"
)

type NewExchangerDeliver struct {
	Redis DomainLevel.ReadingRedis
}

type NewExchangerParsers struct {
	Decoder RepoParsers.Decode
}
type NewExchangerCrypto struct {
	Decrypter  DomainLevel.Decryption
	Validation DomainLevel.CryptoValidating
}
type NewExchangerKeys struct {
	ServerKeys    DomainLevel.NewSetKeys
	EncrypterKeys RepoEncrypterKeys.Keys
}
type NewExchanger struct {
	NewExchangerDeliver
	NewExchangerParsers
	NewExchangerCrypto
	NewExchangerKeys
}

func GetNewExchanger(newExchangerDeliver NewExchangerDeliver, newExchangerParsers NewExchangerParsers, newExchangerCrypto NewExchangerCrypto, newExchangerKeys NewExchangerKeys) *NewExchanger {
	return &NewExchanger{NewExchangerDeliver: newExchangerDeliver, NewExchangerParsers: newExchangerParsers, NewExchangerCrypto: newExchangerCrypto, NewExchangerKeys: newExchangerKeys}
}

type Exchanger interface {
	setCheckedData([]byte) *NewExchangerPacketDetailsOutData
	GetPlanningExchanger() time.Duration
}
type NewExchangerPacketDetailsOutData struct {
	Time   time.Duration
	Error  error
	NewKey *memguard.LockedBuffer
}

func (n NewExchanger) setCheckedData(bytes []byte) *NewExchangerPacketDetailsOutData {
	grpcPacket := &Dto.RedisPacketStructFromMasterServer{
		AesKey:          nil,
		PlainText:       nil,
		Signature:       nil,
		TimeNextSwaping: time.Duration(0),
	}

	err := n.Decoder.JsonDecodeMarshall(&grpcPacket, bytes)
	if err != nil {
		return &NewExchangerPacketDetailsOutData{Error: err}
	}

	AesKeyDecrypted1, err2 := n.Decrypter.DecryptAesKey(n.ServerKeys.GerOurPrivateKey(), grpcPacket.AesKey)
	if err2 != nil {
		return &NewExchangerPacketDetailsOutData{Error: err2}
	}
	NewRsaKey := n.Decrypter.DecryptPacket(AesKeyDecrypted1, grpcPacket.PlainText)
	if NewRsaKey == nil {
		return &NewExchangerPacketDetailsOutData{Error: errors.New(DomainLevel.ErrorDecryptKeys)}
	}
	defer NewRsaKey.Destroy()

	hashSha := sha256.New()
	hashSha.Write(NewRsaKey.Bytes())

	getMasterPublicKey := n.ServerKeys.GetMasterPublicKey()

	err = n.Validation.CheckSignKey(DomainLevel.CheckSignKeyIncomingData{
		Sign:            grpcPacket.Signature,
		Hash:            hashSha.Sum(nil),
		MasterPublicKey: getMasterPublicKey,
	})
	if err != nil {
		return &NewExchangerPacketDetailsOutData{Error: err2}
	}

	return &NewExchangerPacketDetailsOutData{
		Time:   grpcPacket.TimeNextSwaping,
		Error:  nil,
		NewKey: NewRsaKey,
	}

}
func (n NewExchanger) getTimeOutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 6*time.Second)
}
func (n NewExchanger) GetPlanningExchanger() time.Duration {
	ctx, cancel := n.getTimeOutContext()
	defer cancel()
	outData, err := n.Redis.GetKey(ctx)
	if err != nil {
		return DomainLevel.DefaultErrorTime
	}

	data := n.setCheckedData(outData)
	if data.Error != nil {
		return DomainLevel.DefaultErrorTime
	}
	defer data.NewKey.Destroy()
	n.EncrypterKeys.UpdateOldKey()
	err = n.EncrypterKeys.UpdateNewKey(data.NewKey)
	return data.Time
}
