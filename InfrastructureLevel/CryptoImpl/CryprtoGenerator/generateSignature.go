package CryprtoGenerator

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"log/slog"
)

func (s *CryprtoGenerating) SignerData(HashFromData []byte, privateKeyServer []byte) ([]byte, error) {

	slog.Info("Signer Data")
	MasterServerPrivateKey, err := x509.ParsePKCS1PrivateKey(privateKeyServer)
	if err != nil {
		slog.Error("Error Parsing Private Key Server", "Error", err.Error())
		return nil, err
	}

	signed, err := rsa.SignPKCS1v15(rand.Reader, MasterServerPrivateKey, crypto.SHA256, HashFromData)
	if err != nil {
		slog.Error("Error Signing Private Key Server", "Error", err.Error())
		return nil, err
	}

	return signed, nil
}
