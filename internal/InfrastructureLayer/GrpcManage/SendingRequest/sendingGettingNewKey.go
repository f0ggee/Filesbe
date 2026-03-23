package SendingRequest

import (
	pb "Kaban/internal/InfrastructureLayer/GrpcManage/protoFiles"
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *SenderRequests) RequestingGettingNewKey(data []byte) ([]byte, error) {
	slog.Info("We started sending data")
	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Error while creating gRPC connection", "Error", err)
		return nil, err
	}

	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	clientRequest := pb.NewSendingGettingClient(conn)

	OutputData, err := clientRequest.GettingNewKey(ctx, &pb.InputSendData{SendData: data})
	if err != nil {
		slog.Error("Error while sending data", "Error", err)
		return nil, err
	}

	if OutputData.Error != nil {
		slog.Error("we got the error ", "Error", OutputData.Error)
		return nil, err
	}

	return OutputData.BytesOutput, nil
}
