package GrpcManage

import (
	pb "Kaban/internal/InfrastructureLayer/GrpcManage/protoFiles"
	"context"
	"log/slog"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NewSenderRequests struct{}

func GetNewSenderRequests() *NewSenderRequests {
	return &NewSenderRequests{}
}

func (s NewSenderRequests) RequestingGettingNewKey(data []byte) ([]byte, error) {
	credentials := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(os.Getenv("GRPC_ADDR"), credentials)
	if err != nil {
		slog.Error("Error while creating gRPC connection", "Error", err)
		return nil, err
	}
	defer conn.Close()

	clientRequest := pb.NewSendingGettingClient(conn)

	OutputData, err := clientRequest.GetNewKey(context.Background(), &pb.InputSendData{SendData: data})
	if err != nil {
		slog.Error("Error while sending data", "Error", err)
		return nil, err
	}
	if OutputData.Error != nil {
		slog.Error("Got the error", "Error", OutputData.Error)
		return nil, err
	}
	return OutputData.BytesOutput, nil
}
