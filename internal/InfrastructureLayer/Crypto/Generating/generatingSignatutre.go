package Generating

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"log/slog"
)

func (g Generating) GenerateSignature(message []byte, key []byte) ([]byte, error) {
	slog.Info("Starting generating the signature")
	KeyPrivate, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		slog.Error("Error while converting key to private key", "err", err)
		return nil, err
	}
	HashData := sha256.Sum256(message)
	SignedMessage, err := rsa.SignPKCS1v15(rand.Reader, KeyPrivate, crypto.SHA256, HashData[:])
	if err != nil {
		slog.Error("Error while signing message", "err", err)
		return nil, err
	}
	return SignedMessage, nil

}
