package Dto

import "time"

type RedisPacketStructFromMasterServer struct {
	AesKey          []byte        `redis:"aes"`
	PlainText       []byte        `redis:"plaintext"`
	Signature       []byte        `redis:"signature"`
	TimeNextSwaping time.Duration `redis:"time_next_swaping"`
}
