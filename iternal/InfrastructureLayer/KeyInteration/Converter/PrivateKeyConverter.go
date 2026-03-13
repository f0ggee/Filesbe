package Converter

import (
	"crypto/rsa"
	"crypto/x509"
)

type KeyConverter struct{}

func (k KeyConverter) ConverterToPublicKey(Public []byte) (*rsa.PublicKey, error) {
	//TODO implement me
	return x509.ParsePKCS1PublicKey(Public)
}

func (k KeyConverter) ConverterToPrivateKey(bytes []byte) (*rsa.PrivateKey, error) {
	return x509.ParsePKCS1PrivateKey(bytes)
}
