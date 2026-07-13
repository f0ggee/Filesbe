package Dto

import "time"

type GrpcIncomingPacketDetails struct {
	Sign    []byte        `json:"Sign"`
	RsaKey  []byte        `json:"MasterPublicKey"`
	T1      time.Duration `json:"T1"`
	TimeNow time.Time     `json:"TimeNow"`
}

func GetGrpcIncomingPacketDetails() *GrpcIncomingPacketDetails {

	return &GrpcIncomingPacketDetails{
		Sign:    nil,
		RsaKey:  nil,
		T1:      0,
		TimeNow: time.Now(),
	}
}
