package Handlers

import (
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

func (sa *HandlerPackCollect) SwapKeyFirst() time.Duration {

	slog.Info("SwapKeyFirst", "start", true)
	SignedServerName, err := sa.Crypto.Generate.GenerateSignature([]byte(os.Getenv("serverName")), ControlPrivateKeyStruct.OurPrivateKeyIntoBytes)
	if err != nil {
		return 0
	}
	GrpcStruct := Dto.GrpcOutComingPacketDetails{
		Time:             time.Now(),
		SignedServerName: SignedServerName,
		ServerName:       []byte(os.Getenv("serverName")),
	}

	AesKey, err := memguard.NewBufferFromReader(rand.Reader, 32)
	if err != nil {
		slog.Error("Error while generating AesKey", "err", err)
	}
	defer AesKey.Destroy()

	ConvertedData, err := sa.Convert.Converting.JsonConverter(GrpcStruct)
	if err != nil {
		slog.Error("Error while converting", "err", err)
		return DefaultErrorTime
	}

	EncryptedData := []byte(nil)
	EncryptedDataAesKey := []byte(nil)

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			EncryptedData1, err1 := sa.Crypto.Encrypt.EncryptAes(AesKey.Data(), ConvertedData)
			if err1 != nil {
				slog.Error("Error while encrypt", "err", err1)

				return ctx.Err()
			}
			EncryptedData = EncryptedData1
			return nil
		}

	})
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			Key, err1 := x509.ParsePKCS1PublicKey(ControlPrivateKeyStruct.MasterServerPublicKeyBytes)
			if err1 != nil {
				slog.Error("Error while parsing AesKey", "err", err)
				return ctx.Err()
			}

			EncryptedDataAesKey1, err1 := sa.Crypto.Encrypt.EncryptFileInfo(AesKey.Data(), Key)
			if err1 != nil {
				slog.Error("Error while encrypting Info", "err", err1)
				return err1
			}
			EncryptedDataAesKey = EncryptedDataAesKey1
			return nil
		}

	})

	if err := g.Wait(); err != nil {
		slog.Error("Error while generating AesKey", "err", err)
		return DefaultErrorTime
	}

	convertedDataGrpcDataLooks, err := sa.Convert.Converting.JsonConverter(Dto.GrpcOutComingPacketForSending{
		AesKeyData: EncryptedDataAesKey,
		CipherData: EncryptedData,
	})
	if err != nil {
		return DefaultErrorTime
	}

	attempts, sec := 1, 1

	for {
		if attempts >= 12 {
			return 12 * time.Hour
		}
		OutputData, err := sa.Grpc.GrpcSendingRequest.RequestingGettingNewKey(convertedDataGrpcDataLooks)
		if err != nil {
			slog.Error("Error while SendRequestGrpc", "err", err)
			return DefaultErrorTime
		}
		TimeNextSwapping, err := sa.Grpc.ProcessingRequests.CheckingGettingNewKey(OutputData)
		if err != nil {
			attempts++
			sec++
			time.Sleep(time.Duration(sec) + time.Second)
			continue
		}
		if TimeNextSwapping == 0 {
			return DefaultErrorTime
		}
		return TimeNextSwapping
	}
}
