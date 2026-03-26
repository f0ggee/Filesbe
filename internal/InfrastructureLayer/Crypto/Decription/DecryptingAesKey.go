package Decription

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"log/slog"
)

func (d DecryptionData) DecryptAesKey(RsaKey []byte, aesKey []byte) ([]byte, error) {
	slog.Info("DecryptAesKey", "id", rand.Text())
	RsaKeyPrivate, err := x509.ParsePKCS1PrivateKey(RsaKey)
	if err != nil {
		slog.Error("Error Parsing RsaKey", "Func: DecryptAesKey", "Error", err)
		return nil, err
	}

	return rsa.DecryptOAEP(sha256.New(), rand.Reader, RsaKeyPrivate, aesKey, nil)

}
