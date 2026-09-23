package Crypto

import (
	"Kaban/internal/DomainLevel"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"log/slog"
)

const (
	ErrorDecryptFileInfo = "an error happened during decrypting info"
	ErrorDecryptKeys     = "data can't be decrypted because of an unexpected error"
)

type AESDecrypt struct{}

func (A AESDecrypt) DecryptData(data DomainLevel.IncomeData) ([]byte, error) {
	aesBlock, err := aes.NewCipher(data.Key)
	if err != nil {
		slog.Error("Func DecryptPacket:Error create new aes block", "Error", err.Error())
		return nil, errors.New(ErrorDecryptFileInfo)
	}
	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		slog.Error("Func DecryptPacket: Error create new gcm", "Error", err.Error())
		return nil, errors.New(ErrorDecryptFileInfo)
	}
	return gcm.Open(nil, data.Data[:gcm.NonceSize()], data.Data[gcm.NonceSize():], nil)
}

type RsaDecryptOAEP struct{}

func (r RsaDecryptOAEP) DecryptData(data DomainLevel.IncomeData) ([]byte, error) {
	RsaKeyPrivate, err := x509.ParsePKCS1PrivateKey(data.Key)
	if err != nil {
		slog.Error("DecryptAesKey;Error Parsing RsaKey", "Func decrypt error", err)
		return nil, errors.New(ErrorDecryptKeys)
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, RsaKeyPrivate, data.Data, nil)
}
