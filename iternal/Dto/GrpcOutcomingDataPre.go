package Dto

import "time"

type GrpcOutcomingDataPre struct {
	Sign   []byte        `json:"Sign"`
	RsaKey []byte        `json:"RsaKey"`
	T1     time.Duration `json:"T1"`
}
