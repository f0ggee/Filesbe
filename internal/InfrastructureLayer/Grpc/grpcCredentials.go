package Grpc

import (
	"errors"
	"log/slog"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type hub struct {
	g *grpc.ClientConn
}

var connHub hub

func GetConn() *grpc.ClientConn {
	return connHub.g
}

func SetHub(g *grpc.ClientConn) {
	connHub = hub{
		g: g,
	}
}

var ErrorMakeClient = errors.New("error to make a new key request")

func MakeGrpcConn() (*grpc.ClientConn, error) {
	credentials := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(os.Getenv("GRPC_ADDR"), credentials)
	if err != nil {
		slog.Error("MakeNewKeyRequest: error to make a new client", "ERROR", err)
		return nil, ErrorMakeClient
	}
	return conn, nil
}
