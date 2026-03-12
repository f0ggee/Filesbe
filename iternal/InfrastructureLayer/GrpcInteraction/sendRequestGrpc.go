package GrpcInteraction

import (
	"context"
	"time"

	pb "Kaban/iternal/InfrastructureLayer/GrpcInteraction/protoFiles"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DataSend struct {
}

func (s *DataSend) SayHi() string {
	//TODO implement me

	return "hello"
}

func (s *DataSend) SendRequestGrpc(data []byte) ([]byte, error) {
	slog.Info("We've started getting data")
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Error while creating gRPC connection", "Error", err)
		return nil, err
	}

	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	clientRequest := pb.NewSendingGettingClient(conn)

	OutputData, err := clientRequest.GettingNewKey(ctx, &pb.InputSendData{SendData: data})
	if err != nil {
		slog.Error("Error while sending data", "Error", err)
		return nil, err
	}

	if OutputData.Error != nil {
		slog.Error("Error while sending data", "Error", OutputData.Error)
		return nil, err
	}

	return OutputData.BytesOutput, nil
}
