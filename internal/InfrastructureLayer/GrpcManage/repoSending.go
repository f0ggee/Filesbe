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

type NewSenderRequests struct {
	Conn *grpc.ClientConn
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
func (s NewSenderRequests) SetMakerRequestEncrypterKey(convertedDataGrpcDataLooks []byte) ([]byte, error) {
	attempts, sec := 1, 1
	for {
		if attempts > 12 {
			slog.Error("SetMakerRequestEncrypterKey; attempts are expired")
			return nil, errors.New(DomainLevel.ErrorAttemptsExpired)
		}
		OutputData, err := s.SetEncrypterKeyRequest(convertedDataGrpcDataLooks)
		if err != nil {
			slog.Error("SetMakerRequestEncrypterKey; error to send a request. Send another request")
			attempts++
			sec++
			time.Sleep(time.Duration(sec) * time.Second)
			continue
		}
		slog.Info("SetMakerRequestEncrypterKey; data was gotten")
		return OutputData, nil
	}
}
