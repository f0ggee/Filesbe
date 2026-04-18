package Sender

import (
	pb "Kaban/internal/InfrastructureLayer/GrpcManage/protoFiles"
	"context"
	"log/slog"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s SenderRequests) RequestingGettingNewKey(data []byte) ([]byte, error) {
	slog.Info("Start a request for a key")

	conn, err := grpc.NewClient(os.Getenv("GRPC_ADDR"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Error while creating gRPC connection", "Error", err)
		return nil, err
	}
	defer conn.Close()

	//ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//defer cancel()
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
