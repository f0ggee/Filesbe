package Dto

import "time"

type GrpcOutcomingDataPre struct {
	Sign   []byte `json:"Sign"`
	RsaKey []byte `json:"RsaKey"`
	T1     time.Time
	T2     time.Time
	T3     time.Time
}
