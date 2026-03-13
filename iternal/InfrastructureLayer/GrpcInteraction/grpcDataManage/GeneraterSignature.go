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

func (d *DataManage) GenerateSignature(message []byte, key []byte) ([]byte, error) {

	slog.Info("We start the generateSignature")

	KeyPrivate, err := d.K.Keys.ConverterToPrivateKey(key)
	if err != nil {
		slog.Error("Error while converting key to private key", "err", err)
		return nil, err
	}

	HashData := d.K.Keys.ConvertDataToHash(message, nil)

	SignedMessage, err := rsa.SignPKCS1v15(rand.Reader, KeyPrivate, crypto.SHA256, HashData)
	if err != nil {
		slog.Error("Error while signing message", "err", err)
		return nil, err
	}
	return SignedMessage, nil

}
