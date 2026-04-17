package GrpcConn

import (
	"MasterServer_/InfrastructureLevel/CryptoImpl/CryprtoGenerator"
	"MasterServer_/InfrastructureLevel/CryptoImpl/Decryptor"
	"MasterServer_/InfrastructureLevel/CryptoImpl/Encrypter"
	"MasterServer_/InfrastructureLevel/Grpc/GrpcHandleData"
	pbRealization "MasterServer_/InfrastructureLevel/Grpc/HandlerRequests"
	"MasterServer_/InfrastructureLevel/Grpc/PacketValidation"
	pbProtoFiles "MasterServer_/InfrastructureLevel/Grpc/Proto/protoFiles"
	"MasterServer_/InfrastructureLevel/serveManage/ConverterData"
	"MasterServer_/InfrastructureLevel/serveManage/GettingInfo"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

func (g GrpcConnection) GrpcHandleRequests(GrpcHandlingData GrpcHandleData.GrpcDataManagement, ServerInfo GettingInfo.SeverManage, Encrypting Encrypter.Encryption, PacketValidating PacketValidation.ValidatePacketData, CryptoGenerate CryprtoGenerator.CryprtoGenerating, ConvertData ConverterData.ConvertingData, Decrypted Decryptor.Decrypting) {
	lis, err := net.Listen(g.Network, g.Address)
	if err != nil {
		slog.Error("failed to listen:", err.Error())
		return
	}
	grpcServer := grpc.NewServer()
	pbProtoFiles.RegisterSendingGettingServer(grpcServer, &pbRealization.GrpcHandlerGettingNewKey{
		UnimplementedSendingGettingServer: pbProtoFiles.UnimplementedSendingGettingServer{},
		S: pbRealization.HandlingRequestsForNewKey{
			Grpc:             &GrpcHandlingData,
			ServerManagement: &ServerInfo,
			Encryption:       &Encrypting,
			Checker:          &PacketValidating,
			CryptoGenerating: &CryptoGenerate,
			ConverterJson:    &ConvertData,
			Decrypting:       &Decrypted,
		},
	})
	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("failed to serve:", err.Error())
		return
	}

}
