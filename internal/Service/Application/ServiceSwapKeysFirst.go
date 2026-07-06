package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"context"
	"crypto/rand"
	"crypto/x509"
	"log/slog"
	"os"
	"time"

	"github.com/awnumar/memguard"
	"golang.org/x/sync/errgroup"
)

type NewSwapKeyFirst struct {
	GetCrypto
	GetControlKeys
	Parser
}

func GetNewNewSwapKeyFirst(crypto GetCrypto, controlKeys GetControlKeys, parser Parser) *NewSwapKeyFirst {
	return &NewSwapKeyFirst{GetCrypto: crypto, GetControlKeys: controlKeys, Parser: parser}
}

func (sa *NewSwapKeyFirst) GetSwapKeyFirst() time.Duration {

	slog.Info("Func GetSwapKeyFirst:", "start", true)
	serverName := []byte(os.Getenv("serverName"))
	key, err := sa.GetControlKeys.Keys.GerOurPrivateKey()
	if err != nil {
		return DomainLevel.DefaultErrorTime
	}
	SignedServerName, err := sa.GetCrypto.Generate.GenerateSignature(serverName, key)
	if err != nil {
		return 0
	}
	GrpcStruct := Dto.GrpcOutComingPacketDetails{
		Time:             time.Now(),
		SignedServerName: SignedServerName,
		ServerName:       serverName,
	}

	AesKey, err := memguard.NewBufferFromReader(rand.Reader, 32)
	if err != nil {
		slog.Error("Error while generating AesKey", "err", err)
	}
	defer AesKey.Destroy()

	ConvertedData, err := sa.Encode.JsonEncodeMarshall(GrpcStruct)
	if err != nil {
		slog.Error("Error while converting", "err", err)
		return DomainLevel.DefaultErrorTime
	}
	EncryptedData := []byte(nil)
	EncryptedDataAesKey := []byte(nil)
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			EncryptedData1, err1 := sa.GetCrypto.Encrypt.EncryptAes(AesKey.Data(), ConvertedData)
			if err1 != nil {
				slog.Error("Error while encrypt", "err", err1)

				return ctx.Err()
			}
			Sa := &EncryptedData1
			EncryptedData = *Sa
			return nil
		}

	})
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			masterPublicKey, err := sa.GetControlKeys.Keys.GetMasterPublicKey()
			if err != nil {
				return err
			}
			Key, err1 := x509.ParsePKCS1PublicKey(masterPublicKey)
			if err1 != nil {
				slog.Error("Error while parsing Master Server's public masterPublicKey", "err", err)
				return err1
			}

			EncryptedDataAesKey1, err2 := sa.GetCrypto.Encrypt.EncryptFileInfo(AesKey.Data(), Key)
			if err2 != nil {
				slog.Error("Error while encrypting Info", "err", err1)
				return err1
			}

			EncryptedDataAesKey = EncryptedDataAesKey1
			return nil
		}

	})

	if err := g.Wait(); err != nil {
		slog.Error("Error while generating AesKey", "err", err)
		return DomainLevel.DefaultErrorTime
	}

	if EncryptedData == nil || EncryptedDataAesKey == nil {
		slog.Error("Error with data", slog.Group("Data", slog.Any("AesKey", EncryptedDataAesKey), slog.Any("EncryptedData", EncryptedData)))
		return DomainLevel.DefaultErrorTime
	}

	convertedDataGrpcDataLooks, err := sa.Parser.Encode.JsonEncodeMarshall(Dto.GrpcOutComingPacketForSending{
		AesKeyData: EncryptedDataAesKey,
		CipherData: EncryptedData,
	})
	if err != nil {
		return DomainLevel.DefaultErrorTime
	}

	return MakerRequests(sa, convertedDataGrpcDataLooks)
}

func (sa *NewSwapKeyFirst) MakerRequests(convertedDataGrpcDataLooks []byte) time.Duration {
	attempts, sec := 1, 1
	for {
		if attempts > 12 {
			return 12 * time.Hour
		}
		OutputData, err := sa.Grpc.GrpcSendingRequest.RequestingGettingNewKey(convertedDataGrpcDataLooks)
		if err != nil {
			slog.Error("Error while SendRequestGrpc", "err", err)
			attempts++
			sec++
			time.Sleep(time.Duration(sec) * time.Second)
			continue
		}
		TimeNextSwapping, err := sa.Grpc.ProcessingRequests.CheckingGettingNewKey(OutputData)
		if err != nil {
			attempts++
			sec++
			time.Sleep(time.Duration(sec) + time.Second)
			continue
		}
		if TimeNextSwapping == 0 {
			return DomainLevel.DefaultErrorTime
		}
		return TimeNextSwapping
	}
}
