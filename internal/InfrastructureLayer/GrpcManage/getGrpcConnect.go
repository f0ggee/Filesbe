package GrpcManage

import (
	"log/slog"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetNewGrpcConnect() *grpc.ClientConn {
	credentials := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(os.Getenv("GRPC_ADDR"), credentials)
	if err != nil {
		slog.Error("GetNewGrpcConnect; error to create a connect", "ERROR", err)
		return nil
	}
	return conn
}
