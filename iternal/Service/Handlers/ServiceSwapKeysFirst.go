package Handlers

import (
	"Kaban/iternal/Dto"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/awnumar/memguard"
)

func (sa *HandlerPackCollect) SwapKeyFirst() {

	slog.Info("SwapKeyFirst", "start", true)

	OurKeyInPrivateType, err := x509.ParsePKCS1PrivateKey(ControlPrivateKeyStruct.OurPrivateKeyIntoBytes)
	if err != nil {
		slog.Error("Error while parsing OurPrivateKeyIntoBytes", "err", err)
		return
	}

	SignedServerName, err := rsa.SignPKCS1v15(rand.Reader, OurKeyInPrivateType, crypto.SHA256, sha256.New().Sum([]byte(os.Getenv("serverName"))))
	if err != nil {
		slog.Error("Error while signing OurPrivateKeyIntoBytes", "err", err)
		return
	}

	GrpcStruct := Dto.GrpcDataLook{
		Time:             time.Now(),
		SignedServerName: SignedServerName,
		ServerName:       []byte(os.Getenv("serverName")),
	}

	AesKey, err := memguard.NewBufferFromReader(rand.Reader, 32)
	if err != nil {
		slog.Error("Error while generating AesKey", "err", err)
	}
	defer AesKey.Destroy()

	ConvertedData, err := sa.S.Choose.JsonConverter(GrpcStruct)
	if err != nil {
		return
	}

	EncryptedData := []byte(nil)
	EncryptedDataAesKey := []byte(nil)

	chanelForErrors := make(chan error, 2)

	var wg sync.WaitGroup

	for i := 1; i <= 1; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			EncryptedData1, err1 := sa.S.Encryption.EncryptAes(AesKey.Data(), ConvertedData)
			if err1 != nil {
				chanelForErrors <- err1
				return
			}
			EncryptedData = EncryptedData1
		}()
		go func() {
			defer wg.Done()
			EncryptedDataAesKey1, err1 := sa.S.FileDataManipulation.EncryptData(AesKey.Data(), ControlPrivateKeyStruct.MasterServerPublicKeyBytes)
			if err1 != nil {
				chanelForErrors <- err1
				return
			}
			EncryptedDataAesKey = EncryptedDataAesKey1
		}()
	}

	wg.Wait()
	close(chanelForErrors)
	for err := range chanelForErrors {
		if err != nil {
			return
		}
	}

	convertedDataGrpcDataLooks, err := sa.S.Choose.JsonConverter(Dto.GrpcSendingData{
		AesKeyData: EncryptedDataAesKey,
		CipherData: EncryptedData,
	})
	if err != nil {
		return
	}

}
