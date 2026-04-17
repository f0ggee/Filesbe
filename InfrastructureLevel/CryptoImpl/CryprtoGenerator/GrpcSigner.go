package CryprtoGenerator

import (
	"MasterServer_/Dto"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"log/slog"
)

func (s *CryprtoGenerating) GrpcSignerKey() ([]byte, error) {

	privateKeyBytes, err := hex.DecodeString(Dto.Keys.MasterServerKey)
	if err != nil {
		slog.Error("GenerateSignatureFromKey Error while trying to decode private key", "error", err.Error())
		return nil, err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		slog.Error("Error while trying to parse private key", "Func:GenerateSignature", "error", err.Error())
		return nil, err
	}

	hash := sha256.Sum256(Dto.Keys.NewPrivateKey.Bytes())
	SignedMassage, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		slog.Error("GenerateSignatureFromKey Error while trying to sign signature", "error", err.Error())
		return nil, err
	}
	return SignedMassage, nil

}
