package ProtocolManage

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/awnumar/memguard"
	"golang.org/x/sync/errgroup"
)

type NewExchangeInitializerCrypto struct {
	CryptoGenerating DomainLevel.CryptoGenerating
	CryptoEncrypt    DomainLevel.Encryption
	CryptoDecrypt    DomainLevel.Decryption
	CryptoValidate   DomainLevel.CryptoValidating
}

func GetNewExchangeInitializerCrypto(cryptoGenerating DomainLevel.CryptoGenerating, cryptoEncrypt DomainLevel.Encryption, cryptoDecrypt DomainLevel.Decryption, cryptoValidate DomainLevel.CryptoValidating) *NewExchangeInitializerCrypto {
	return &NewExchangeInitializerCrypto{CryptoGenerating: cryptoGenerating, CryptoEncrypt: cryptoEncrypt, CryptoDecrypt: cryptoDecrypt, CryptoValidate: cryptoValidate}
}

type NewExchangeInitializerParsers struct {
	Encode RepoParsers.Encode
	Decode RepoParsers.Decode
}

func GetNewExchangeInitializerParsers(encode RepoParsers.Encode, decode RepoParsers.Decode) *NewExchangeInitializerParsers {
	return &NewExchangeInitializerParsers{Encode: encode, Decode: decode}
}

type NewExchangeInitializerDeliver struct {
	Grcp DomainLevel.Requests
}

func GetNewExchangeInitializerDeliver(grcp DomainLevel.Requests) *NewExchangeInitializerDeliver {
	return &NewExchangeInitializerDeliver{Grcp: grcp}
}

type NewExchangeInitializerKey struct {
	Keys       RepoEncrypterKeys.Keys
	ServerKeys DomainLevel.NewServerKeys
}

func GetNewExchangeInitializerKey(keys RepoEncrypterKeys.Keys, serverKeys DomainLevel.NewServerKeys) *NewExchangeInitializerKey {
	return &NewExchangeInitializerKey{Keys: keys, ServerKeys: serverKeys}
}

type NewExchangeInitializer struct {
	NewExchangeInitializerKey
	NewExchangeInitializerCrypto
	NewExchangeInitializerParsers
	NewExchangeInitializerDeliver
}

func GetNewExchangeInitializer(newExchangeInitializerKey NewExchangeInitializerKey, newExchangeInitializerCrypto NewExchangeInitializerCrypto, newExchangeInitializerParsers NewExchangeInitializerParsers, newExchangeInitializerConnect NewExchangeInitializerDeliver) *NewExchangeInitializer {
	return &NewExchangeInitializer{NewExchangeInitializerKey: newExchangeInitializerKey, NewExchangeInitializerCrypto: newExchangeInitializerCrypto, NewExchangeInitializerParsers: newExchangeInitializerParsers, NewExchangeInitializerDeliver: newExchangeInitializerConnect}
}

type PacketDetailsOutData struct {
	Time   time.Duration
	Error  error
	NewKey *memguard.LockedBuffer
}
type ExchangeInitializer interface {
	getPreparingData() ([]byte, error)
	getPacketDetails([]byte) *PacketDetailsOutData
	GetExchangerInitializer() time.Duration
}

func (n *NewExchangeInitializer) GetExchangerInitializer() time.Duration {
	Data, err := n.getPreparingData()
	if err != nil {
		return DefaultErrorTime
	}
	outData, err := n.Grcp.SetEncrypterKeyRequest(Data)
	if err != nil {
		return DefaultErrorTime
	}
	PacketDetails := n.getPacketDetails(outData)
	if PacketDetails.Error != nil {
		return DefaultErrorTime
	}
	n.Keys.UpdateOldKey()
	n.Keys.UpdateNewKey(PacketDetails.NewKey)
	PacketDetails.NewKey.Destroy()
	return PacketDetails.Time
}

func (n *NewExchangeInitializer) getPreparingData() ([]byte, error) {
	serverName := []byte(os.Getenv("serverName"))
	key := n.ServerKeys.GerOurPrivateKey()
	SignedServerName, err := n.CryptoGenerating.GenerateSignature(serverName, key)
	if err != nil {
		return nil, err
	}
	GrpcStruct := Dto.GetGrpcOutComingPacketDetails(time.Now(), serverName, SignedServerName)
	AesKey, err := memguard.NewBufferFromReader(rand.Reader, 32)
	if err != nil {
		slog.Error("getPreparingData; error to generate an AES-KEY", "ERROR", err)
		return nil, err
	}
	defer AesKey.Destroy()

	ConvertedData, err := n.Encode.JsonEncodeMarshall(GrpcStruct)
	if err != nil {
		return nil, err
	}

	EncryptedData := []byte(nil)
	EncryptedDataAesKey := []byte(nil)
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			EncryptedData1, err1 := n.CryptoEncrypt.EncryptAes(AesKey.Data(), ConvertedData)
			if err1 != nil {
				return ctx.Err()
			}
			EncryptedData = *&EncryptedData1
			return nil
		}

	})
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			masterPublicKey := n.ServerKeys.GetMasterPublicKey()
			Key, err1 := x509.ParsePKCS1PublicKey(masterPublicKey)
			if err1 != nil {
				slog.Error("Error while parsing Master Server's public masterPublicKey", "err", err)
				return err1
			}

			EncryptedDataAesKey1, err2 := n.CryptoEncrypt.EncryptFileInfo(AesKey.Data(), Key)
			if err2 != nil {
				return err1
			}
			EncryptedDataAesKey = *&EncryptedDataAesKey1
			return nil
		}

	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if EncryptedData == nil || EncryptedDataAesKey == nil {
		slog.Error("getPreparingData; error to get data")
		return nil, errors.New(DomainLevel.ErrorDataNil)
	}

	return n.Encode.JsonEncodeMarshall(Dto.GrpcOutComingPacketForSending{
		AesKeyData: EncryptedDataAesKey,
		CipherData: EncryptedData,
	})
}

func (n NewExchangeInitializer) getPacketDetails(bytes []byte) *PacketDetailsOutData {
	Id := rand.Text()
	slog.Info("Func CheckingGettingNewKey: Start accepting the new key", slog.String("ID", Id))
	PacketLook := Dto.GetNewGrpcOutComingPacketForSending()

	err := n.Decode.JsonDecodeMarshall(&PacketLook, bytes)
	if err != nil {
		return &PacketDetailsOutData{Error: err}
	}

	DecryptedAesKey, err := n.CryptoDecrypt.DecryptAesKey(n.ServerKeys.GerOurPrivateKey(), PacketLook.AesKeyData)
	if err != nil {
		return &PacketDetailsOutData{Error: err}
	}

	PacketData := n.CryptoDecrypt.DecryptPacket(DecryptedAesKey, PacketLook.CipherData)
	if PacketData == nil {
		return &PacketDetailsOutData{Error: err}
	}
	defer PacketData.Destroy()

	PacketInfo := Dto.GetGrpcIncomingPacketDetails()
	err = n.Decode.JsonDecodeMarshall(&PacketInfo, PacketData.Bytes())
	if err != nil {
		return &PacketDetailsOutData{Error: err}
	}
	packetTime := PacketInfo.T1
	err = n.setCheckTime(PacketInfo.TimeNow)
	if err != nil {
		return &PacketDetailsOutData{Error: err}
	}

	NewSavingRsa, err := n.getSaveKey(&PacketInfo.RsaKey)
	if err != nil {
		return &PacketDetailsOutData{Error: err}
	}

	Hash := sha256.New()
	Hash.Write(NewSavingRsa.Bytes())
	err = n.CryptoValidate.CheckSignKey(DomainLevel.CheckSignKeyIncomingData{
		Sign:            PacketInfo.Sign,
		Hash:            Hash.Sum(nil),
		MasterPublicKey: n.ServerKeys.GetMasterPublicKey(),
	})
	if err != nil {
		return &PacketDetailsOutData{Error: err}
	}

	return &PacketDetailsOutData{
		Time:   packetTime,
		Error:  nil,
		NewKey: NewSavingRsa,
	}
}

func (n *NewExchangeInitializer) setCheckTime(TimePacket time.Time) error {
	if time.Now().Before(TimePacket) {
		return errors.New(ErrorTimePacket)
	}
	return nil
}
func (h NewExchangeInitializer) getSaveKey(data *[]byte) (*memguard.LockedBuffer, error) {
	if data == nil {
		return nil, errors.New(DomainLevel.ErrorDataNil)
	}
	NewSavingRsa := memguard.NewBuffer(len(*data))
	NewSavingRsa.Copy(*data)
	memguard.WipeBytes(*data)
	return NewSavingRsa, nil
}
