package Encrypter

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"log/slog"
)

func (s *Encryption) EncryptAesKey(AesKey []byte, RsaKey []byte) ([]byte, error) {
	RsaKeyInPublic, err := x509.ParsePKCS1PublicKey(RsaKey)
	if err != nil {
		slog.Error("Error Parsing RsaKey", "Func: EncryptAesKey", "Error", err.Error())
		return nil, err
	}
	encryptedText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, RsaKeyInPublic, AesKey, nil)
	if err != nil {
		slog.Error("Error Encrypting", "Error", err.Error())
		return nil, err
	}
	return encryptedText, nil
}
