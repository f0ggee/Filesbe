package GlobalProces

import (
	"MasterServer_/DomainLevel"
	InftarctionLevel "MasterServer_/InfrastructureLevel"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/awnumar/memguard"
)

func PacketPack(EncryptedAesKey []byte, EncryptedRsaKey []byte, Sign []byte, TimeNextSwaping time.Duration) ([]byte, error) {
	RedisDat := &DomainLevel.RedisDataLooksLike{
		AesKey:          EncryptedAesKey,
		PlainText:       EncryptedRsaKey,
		Signature:       Sign,
		TimeNextSwaping: TimeNextSwaping,
	}

	DataJson, err := json.Marshal(RedisDat)
	if err != nil {
		slog.Error("Error marshalling RedisDat", "Error", err.Error())
		return nil, err
	}
	return DataJson, nil
}

func (psa *ControllingExchange) SwapKeys(KeyServer []byte, RsaKeyNew []byte, NameServer string) error {

	AesKey, err := memguard.NewBufferFromReader(rand.Reader, 32)
	if err != nil {
		slog.Error("Error creating the AesKey", "Error", err.Error())
		return err
	}
	defer AesKey.Destroy()
	HashSha := sha256.New()
	HashSha.Write(RsaKeyNew)

	Sign, err := psa.E.CryptoGen.SignerData(HashSha.Sum([]byte(nil)), psa.E.CryptoKey.GetMasterKey())
	if err != nil {
		return err
	}

	EncryptedRsaKey, err := psa.E.Cryptos.EncryptRsaKey(AesKey.Bytes(), RsaKeyNew)
	if err != nil {

		return err
	}

	EncryptedAesKey, err := psa.E.Cryptos.EncryptAesKey(AesKey.Bytes(), KeyServer)
	if err != nil {
		return err
	}

	DataJson, err2 := PacketPack(EncryptedAesKey, EncryptedRsaKey, Sign, psa.GetDuration())
	if err2 != nil {
		return err2
	}

	err = psa.E.RedisInteracting.SendData(DataJson, NameServer)
	if err != nil {
		return err
	}

	return nil
}

func (psa *ControllingExchange) GetDuration() time.Duration {
	xz := psa.E.TimeData.GetPreviousSwapTime().Add(InftarctionLevel.TimeForSwapping)
	TimeUntil := time.Until(xz)

	return TimeUntil
}
func (psa ControllingExchange) SwapKeysTest(KeyServer []byte, RsaKeyNew []byte, NameServer string) error {
	if NameServer == "" {
		return errors.New("NameServer is required")

	}
	return nil
}
