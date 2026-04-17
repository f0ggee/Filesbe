package Dto

import "time"

type GrpcOutcomingDataPacket struct {
	Sign    []byte        `json:"Sign"`
	RsaKey  []byte        `json:"Key"`
	T2      time.Duration `json:"T1"`
	TimeNow time.Time     `json:"TimeNow"`
}
