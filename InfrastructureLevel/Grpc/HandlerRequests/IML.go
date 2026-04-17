package HandlerRequests

import (
	"MasterServer_/DomainLevel"
	pb "MasterServer_/InfrastructureLevel/Grpc/Proto/protoFiles"
)

type HandlingRequestsForNewKey struct {
	Grpc             DomainLevel.GrpcHandleData
	ServerManagement DomainLevel.GettingServersInfo
	Encryption       DomainLevel.Encryption
	Checker          DomainLevel.PacketChecker
	Decrypting       DomainLevel.Decryptor
	CryptoGenerating DomainLevel.CryptoGenerator
	ConverterJson    DomainLevel.ConverterData
}

type GrpcHandlerGettingNewKey struct {
	pb.UnimplementedSendingGettingServer
	S HandlingRequestsForNewKey
}

func NewGrpcHandlerGettingNewKey(unimplementedSendingGettingServer *pb.UnimplementedSendingGettingServer, s *HandlingRequestsForNewKey) *GrpcHandlerGettingNewKey {
	return &GrpcHandlerGettingNewKey{UnimplementedSendingGettingServer: *unimplementedSendingGettingServer, S: *s}
}
