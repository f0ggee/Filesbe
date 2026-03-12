package Converter

import (
	"crypto/rsa"
	"crypto/x509"
)

type KeyConverter struct{}

func (k KeyConverter) ConverterToPrivateKey(bytes []byte) (*rsa.PrivateKey, error) {
	return x509.ParsePKCS1PrivateKey(bytes)
}
