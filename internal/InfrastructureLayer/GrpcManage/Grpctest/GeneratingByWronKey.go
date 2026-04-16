package Grpctest

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

type EncryptWrongKEy struct {
}

func (e EncryptWrongKEy) EncryptByWrongKey(bytes []byte) ([]byte, error) {

	Key, _ := rsa.GenerateKey(rand.Reader, 2048)

	NewData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &Key.PublicKey, bytes, nil)

	if err != nil {
		return nil, err
	}

	return NewData, nil
}
