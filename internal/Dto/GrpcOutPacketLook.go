package Dto

import "time"

type GrpcOutComingPacketDetails struct {
	Time             time.Time
	ServerName       []byte
	SignedServerName []byte
}

func GetGrpcOutComingPacketDetails(T time.Time, ServerName []byte, SignedName []byte) *GrpcOutComingPacketDetails {

	return &GrpcOutComingPacketDetails{
		Time:             T,
		ServerName:       ServerName,
		SignedServerName: SignedName,
	}
}
