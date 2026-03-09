package Dto

import "time"

type GrpcDataLook struct {
	Time             time.Time
	ServerName       []byte
	SignedServerName []byte
}
