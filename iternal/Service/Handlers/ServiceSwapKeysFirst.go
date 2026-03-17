package Handlers

import (
	"Kaban/iternal/Dto"
	"crypto/rand"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/awnumar/memguard"
)

func (sa *HandlerPackCollect) SwapKeyFirst() time.Duration {

	slog.Info("SwapKeyFirst", "start", true)

	SignedServerName, err := sa.S.GrpcDataMange.GenerateSignature([]byte(os.Getenv("serverName")), ControlPrivateKeyStruct.OurPrivateKeyIntoBytes)
	if err != nil {
		return 0
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

	ConvertedData, err := sa.S.ConverterKey.JsonConverter(GrpcStruct)
	if err != nil {
		return 0
	}

	EncryptedData := []byte(nil)
	EncryptedDataAesKey := []byte(nil)

	chanelForErrors := make(chan error, 2)

	var wg sync.WaitGroup

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
		Key, errConverterKey := sa.S.ConverterKey.ConverterToPublicKey(ControlPrivateKeyStruct.MasterServerPublicKeyBytes)
		if errConverterKey != nil {
			chanelForErrors <- errConverterKey
			return
		}
		EncryptedDataAesKey1, err1 := sa.S.FileDataManipulation.EncryptData(AesKey.Data(), Key)
		if err1 != nil {
			chanelForErrors <- err1
			return
		}
		EncryptedDataAesKey = EncryptedDataAesKey1
	}()

	wg.Wait()
	close(chanelForErrors)
	for err := range chanelForErrors {
		if err != nil {
			return 0
		}
	}

	convertedDataGrpcDataLooks, err := sa.S.ConverterKey.JsonConverter(Dto.GrpcSendingData{
		AesKeyData: EncryptedDataAesKey,
		CipherData: EncryptedData,
	})
	if err != nil {
		return 0
	}

	attempts := 0

	for {
		if attempts > 3 {
			return 12 * time.Hour
		}
		OutputData, err := sa.S.GrpcConn.SendRequestGrpc(convertedDataGrpcDataLooks)
		if err != nil {
			slog.Error("Error while SendRequestGrpc", "err", err)
			return 12 * time.Hour
		}
		TimeForSwapping, err := sa.S.Checking.Handle(OutputData)
		if err != nil {
			attempts++
			time.Sleep(time.Second)
		}
		if TimeForSwapping == 0 {
			return 12 * time.Hour
		}
		return TimeForSwapping
	}
}
