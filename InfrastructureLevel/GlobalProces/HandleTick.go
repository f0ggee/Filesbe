package GlobalProces

import (
	"MasterServer_/Dto"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/awnumar/memguard"
)

func PacketPacking(EncryptedAesKey []byte, EncryptedRsaKey []byte, Sign []byte, err error) ([]byte, error) {
	RedisDat := &Dto.RedisDataLooksLike{
		AesKey:    EncryptedAesKey,
		PlainText: EncryptedRsaKey,
		Signature: Sign,
	}

	DataJson, err := json.Marshal(RedisDat)
	if err != nil {
		slog.Error("Error marshalling RedisDat", "Error", err.Error())
		return nil, err
	}
	return DataJson, nil
}

func (psa *ControllingExchange) SwapKeys(KeyServer []byte, RsaKeyNew []byte, NameServer string) error {

	slog.Info("Starting HandlingRequests")

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

	DataJson, err2 := PacketPacking(EncryptedAesKey, EncryptedRsaKey, Sign, err)
	if err2 != nil {
		return err2
	}

	err = psa.E.RedisInteracting.SendData(DataJson, NameServer)
	if err != nil {
		return err
	}

	return nil
}

func (psa ControllingExchange) SwapKeysTest(KeyServer []byte, RsaKeyNew []byte, NameServer string) error {
	if NameServer == "" {
		return errors.New("NameServer is required")

	}
	return nil
}
