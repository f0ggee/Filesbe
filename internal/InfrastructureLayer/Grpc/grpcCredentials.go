package Grpc

import (
	"errors"
	"log/slog"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var GrpcConn *grpc.ClientConn
var ErrorMakeClient = errors.New("error to make a new key request")

func init() {
	credentials := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(os.Getenv("GRPC_ADDR"), credentials)
	if err != nil {
		slog.Error("MakeNewKeyRequest: error to make a new client", "ERROR", err)
		panic(ErrorMakeClient)
	}
	GrpcConn = conn
}
