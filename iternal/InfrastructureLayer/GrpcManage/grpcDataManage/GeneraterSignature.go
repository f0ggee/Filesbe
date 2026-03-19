package grpcDataManage

import (
	"Kaban/iternal/DomainLevel"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"log/slog"
)

type CollectorPackForGrpcDataManage struct {
	Keys DomainLevel.Converter
}
type DataManage struct {
	K CollectorPackForGrpcDataManage
}

func (d *DataManage) SayHI() string {
	//TODO implement me
	return "Hello"
}
