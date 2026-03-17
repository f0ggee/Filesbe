package GrpcManage

import (
	pb "Kaban/iternal/InfrastructureLayer/GrpcManage/protoFiles"
	"context"
	"time"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DataSend struct{}

func (s *DataSend) SayHi() string {
	//TODO implement me

	return "hello"
}

func (s *DataSend) SendRequestGrpc(data []byte) ([]byte, error) {
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
