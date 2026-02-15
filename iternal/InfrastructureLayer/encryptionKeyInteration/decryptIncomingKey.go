package encryptionKeyInteration

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"log/slog"
)

func (*EncryptionKey) DecryptIncomingKey(aesKey []byte, plainText []byte, ourPrivateKey []byte) []byte {

	privateKey, err := x509.ParsePKCS1PrivateKey(ourPrivateKey)
	if err != nil {
		slog.Error("Error parse ourtPrivateKeyIntoBytes", "Error", err.Error())
		return nil
	}

	AesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, aesKey, nil)
	if err != nil {
		slog.Error("Error decrypt our new private key", "Error", err.Error())
		return nil
	}

	aesBlock, err := aes.NewCipher(AesKey)
	if err != nil {
		slog.Error("Error create new aes block", "Error", err.Error())
		return nil
	}
	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		slog.Error("Error create new gcm", "Error", err.Error())
		return nil
	}

	rsaKey, err := gcm.Open(nil, plainText[:gcm.NonceSize()], plainText[gcm.NonceSize():], nil)
	if err != nil {
		slog.Error("Error decrypt our new private key", "Error", err.Error())
		return nil
	}
	return rsaKey
}
