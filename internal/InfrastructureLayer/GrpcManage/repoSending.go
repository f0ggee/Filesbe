package GrpcManage

import (
	"Kaban/internal/DomainLevel"
	pb "Kaban/internal/InfrastructureLayer/GrpcManage/protoFiles"
	"context"
	"errors"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

var (
	ErrorAttemptsExpired = errors.New("attempts are expired")
)

type NewSenderRequests struct {
	Conn *grpc.ClientConn
}

func (s NewSenderRequests) SetAdditionalData(bytes []byte) DomainLevel.MakerKeyRequest {
	return s
}

func (s NewSenderRequests) Make() (DomainLevel.Requests, error) {
	return s, nil
}

func GetNewSenderRequests() *NewSenderRequests {
	return &NewSenderRequests{}
}

func (s NewSenderRequests) SetEncrypterKeyRequest(data []byte) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientRequest := pb.NewSendingGettingClient(s.Conn)
	OutputData, err := clientRequest.GetNewKey(ctx, &pb.InputSendData{SendData: data})
	if err != nil {
		slog.Error("SetEncrypterKeyRequest; there is an error in the GRPC mechanism", "ERROR", err)
		return nil, err
	}

	return OutputData.BytesOutput, nil
}
func (s NewSenderRequests) SetNewKeyRequest(convertedDataGrpcDataLooks []byte) ([]byte, error) {
	attempts, sec := 1, 1
	for {
		if attempts > 12 {
			slog.Error("SetNewKeyRequest; attempts are expired")
			return nil, ErrorAttemptsExpired
		}
		OutputData, err := s.SetEncrypterKeyRequest(convertedDataGrpcDataLooks)
		if err != nil {
			slog.Error("SetNewKeyRequest; error to send a request. Send another request")
			attempts++
			sec++
			time.Sleep(time.Duration(sec) * time.Second)
			continue
		}
		slog.Info("SetNewKeyRequest; data was gotten")
		return OutputData, nil
	}
}
