package HandlerRequests

import (
	"MasterServer_/DomainLevel"
	"MasterServer_/Dto"
	InftarctionLevel "MasterServer_/InfrastructureLevel"
	pb "MasterServer_/InfrastructureLevel/Grpc/Proto/protoFiles"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/awnumar/memguard"
	"golang.org/x/sync/errgroup"
)

func BreakJsonPacket(err error, Data []byte) (*Dto.GrpcDataIncomingPacket, error) {
	DataIntoPacket := Dto.GrpcDataIncomingPacket{
		Time:             time.Time{},
		ServerName:       nil,
		SignedServerName: nil,
	}

	err = json.Unmarshal(Data, &DataIntoPacket)
	if err != nil {
		slog.Error("Unmarshal Error", "Error", err.Error())
		return nil, errors.New("something gone wrong")
	}
	return &DataIntoPacket, err
}

func DecodePacket(data []byte) (*Dto.GrpcPacket, error) {
	DataIncomingLook := new(Dto.GrpcPacket)
	err := json.Unmarshal(data, &DataIncomingLook)
	if err != nil {
		slog.Error("Unmarshal Error", "Error", err.Error())
		return nil, err
	}
	return DataIncomingLook, nil
}

func (s GrpcHandlerGettingNewKey) GetNewKey(ctx context.Context, data *pb.InputSendData) (*pb.OutputSendData, error) {

	select {
	case <-time.After(time.Second * 10):
		slog.Info("Func GetNewKey: Start exchanging a key")
		if data == nil {
			slog.Error("Data was empty")
			return &pb.OutputSendData{}, errors.New(DomainLevel.DataNil)
		}

		err2 := HashManipulate(s, data)
		if err2 != nil {
			return &pb.OutputSendData{}, err2
		}

		DataIncomingLook, err := DecodePacket(data.SendData)
		if err != nil {
			return &pb.OutputSendData{}, errors.New(DomainLevel.AnotherTypeError)
		}

		DecryptedAesKey, err := s.S.Decrypting.GrpcDecrypterAesKey(DataIncomingLook.AesKeyData)
		if err != nil {
			return &pb.OutputSendData{}, errors.New(DomainLevel.ErrorVerify)
		}

		Data, err := s.S.Decrypting.DecrypterCipherData(DecryptedAesKey, DataIncomingLook.CipherData)
		if err != nil {
			return &pb.OutputSendData{}, errors.New(DomainLevel.ErrorVerify)
		}

		DataIntoPacket, err := BreakJsonPacket(err, Data)
		slog.Info("Func GetNewKey", "The server, which is trying to connect", DataIntoPacket.ServerName)
		if err != nil {
			return &pb.OutputSendData{}, errors.New(DomainLevel.AnotherTypeError)
		}

		ResultComparingTime := s.S.Checker.CheckLifePacket(DataIntoPacket.Time)
		if ResultComparingTime {
			slog.Info("ResultComparingTime", "Time", DataIntoPacket.Time.String())
			return nil, errors.New(DomainLevel.ErrorVerify)
		}

		serversKey := os.Getenv(string(DataIntoPacket.ServerName))
		if serversKey == "" {
			slog.Error("Server Key is empty", "ServerName", string(DataIntoPacket.ServerName))
			return &pb.OutputSendData{}, errors.New(DomainLevel.AnotherTypeError)
		}

		ServerKey, err1 := hex.DecodeString(serversKey)
		if err1 != nil {
			slog.Error("Func GetNewKey: Server Key is invalid", "Error", err1.Error())
			return &pb.OutputSendData{}, errors.New(DomainLevel.AnotherTypeError)
		}

		err = s.S.Checker.CheckSignature(DataIntoPacket.SignedServerName, ServerKey, DataIntoPacket.ServerName)
		if err != nil {
			return &pb.OutputSendData{}, errors.New(DomainLevel.ErrorVerify)
		}

		SignedKey, err := s.S.CryptoGenerating.GrpcSignerKey()
		if err != nil {
			return &pb.OutputSendData{}, errors.New(DomainLevel.AnotherTypeError)
		}

		calculateSwapingTime := s.CalculateSwapingTime()
		Dto.Keys.Mu.Lock()
		OutcomingDataJson, err := s.S.ConverterJson.ConvertDataToJsonType(Dto.GrpcOutcomingDataPacket{
			Sign:   SignedKey,
			RsaKey: Dto.Keys.NewPrivateKey.Bytes(),
			T1:     calculateSwapingTime,
		})
		Dto.Keys.Mu.Unlock()
		if err != nil {
			slog.Error("Marshal Error", "Error", err.Error())
			return &pb.OutputSendData{}, errors.New(DomainLevel.AnotherTypeError)
		}

		AesKey, err := memguard.NewBufferFromReader(rand.Reader, 32)
		if err != nil {
			slog.Error("Generate AesKey Error", "Error", err.Error())
			return &pb.OutputSendData{}, errors.New(DomainLevel.AnotherTypeError)
		}
		defer AesKey.Destroy()

		g, ctx2 := errgroup.WithContext(ctx)

		plainText := []byte{}
		encryptedAesKey := []byte{}

		g.Go(func() error {
			select {
			case <-ctx2.Done():
				return ctx2.Err()

			default:
			}
			PlainText, err12 := s.S.Encryption.EncryptRsaKey(AesKey.Bytes(), OutcomingDataJson)
			if err != nil {
				slog.Error("Func GetNewKey: Encrypt RSA-KEY Error ", "Error", err.Error())
				return err12
			}
			plainText = PlainText
			return nil
		})

		g.Go(func() error {
			select {
			case <-ctx2.Done():
				return ctx2.Err()
			default:
			}
			EncryptedAesKey, err13 := s.S.Encryption.EncryptAesKey(AesKey.Bytes(), ServerKey)
			if err != nil {
				slog.Error("Func GetNewKey: Error encrypt the aesKey", "Error", err)
				return err13
			}
			encryptedAesKey = EncryptedAesKey
			return nil
		})
		if err := g.Wait(); err != nil {
			return &pb.OutputSendData{}, errors.New(DomainLevel.AnotherTypeError)
		}

		OutcomingPacket, err := s.S.ConverterJson.ConvertDataToJsonType(Dto.GrpcPacket{
			AesKeyData: encryptedAesKey,
			CipherData: plainText,
		})
		if err != nil {
			slog.Error("Func GetNewKey: Marshal Error", "Error", err.Error())
			return &pb.OutputSendData{}, errors.New(DomainLevel.ErrorVerify)
		}

		slog.Info("Func GetNewKey: finish an exchange", "The server", DataIntoPacket.ServerName)
		return &pb.OutputSendData{
			BytesOutput: OutcomingPacket,
		}, nil
	}

}

func (S GrpcHandlerGettingNewKey) CalculateSwapingTime() time.Duration {
	slog.Info("Func CalculateSwapingTime", "Time before", S.S.Time.GetPreviousSwapTime())

	xz := S.S.Time.GetPreviousSwapTime().Add(InftarctionLevel.TimeForSwapping)
	s := time.Until(xz)
	slog.Info("Func CalculateSwapingTime", "Time for next swaping", s)
	return s
}
func (S GrpcHandlerGettingNewKey) CalculateSwapingTimeTestFunc(TimeToAdd time.Duration) time.Duration {
	slog.Info("Func CalculateSwapingTime", "Time before", S.S.Time.GetPreviousSwapTime())

	xz := S.S.Time.GetPreviousSwapTime().Add(TimeToAdd)
	s := time.Until(xz)
	slog.Info("Func CalculateSwapingTime", "Time for next swaping", s)
	return s
}

func HashManipulate(s GrpcHandlerGettingNewKey, data *pb.InputSendData) error {
	if s.S.Checker.FindHash([32]byte(sha256.New().Sum(data.SendData[:]))) {
		slog.Error("Hash has been used")
		return errors.New("something gone wrong")
	}

	s.S.Grpc.SaveHash(data.SendData)
	return nil
}
