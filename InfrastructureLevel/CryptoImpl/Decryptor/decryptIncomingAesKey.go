package Decryptor

import (
	"MasterServer_/Dto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"log/slog"
)

func (g Decrypting) GrpcDecrypterAesKey(AesKey []byte) ([]byte, error) {

	KeyInBytes, err := hex.DecodeString(Dto.Keys.MasterServerKey)
	if err != nil {
		slog.Error("Error decoding master server key", err.Error())
		return nil, err
	}

	KeyRsaType, err := x509.ParsePKCS1PrivateKey(KeyInBytes)
	if err != nil {
		slog.Error("Error parsing the master server key", err.Error())
		return nil, err
	}

	DecryptedAesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, KeyRsaType, AesKey, nil)
	if err != nil {
		slog.Error("Error decrypting master server key RSA", err.Error())
		return nil, err
	}
	return DecryptedAesKey, nil
}
