package Dto

import "time"

type GrpcOutComingPacketDetails struct {
	Time             time.Time
	ServerName       []byte
	SignedServerName []byte
}
