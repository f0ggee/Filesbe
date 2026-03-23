package Dto

import "time"

type GrpcIncomingPacketDetails struct {
	Sign   []byte        `json:"Sign"`
	RsaKey []byte        `json:"RsaKey"`
	T1     time.Duration `json:"T1"`
}
