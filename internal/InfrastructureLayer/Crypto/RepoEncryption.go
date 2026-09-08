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
	"fmt"
	"io"
	"log/slog"
	"strings"
)

const ErrorEncryptFile = "can't encrypt a file"

type Encrypter struct {
}

func GetNewEncrypter() Encrypter {
	return Encrypter{}
}

func (*Encrypter) EncryptFileInfo(FileInfoData []byte, Key *rsa.PublicKey) ([]byte, error) {
	encryptAesKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, Key, FileInfoData, nil)
	switch {
	case strings.Contains(fmt.Sprint(err), "message too long for RSA key size"):
		return nil, errors.New(DomainLevel.ErrorFileNameEncrypt)

	case err != nil:
		slog.Error("EncryptFileInfo; error to encrypt", "ERROR", err)
		return nil, errors.New(DomainLevel.ErrorStrangeCrypto)

	}
	return encryptAesKey, nil
}
func (e *Encrypter) EncryptAes(AesKey []byte, Data []byte) ([]byte, error) {
	loggers := slog.With("EncryptAes")
	AesCipher, err := aes.NewCipher(AesKey)
	if err != nil {
		loggers.Error("Error creating new AesCipher", "ERROR", err.Error())
		return nil, err
	}
	NewGcmBlock, err := cipher.NewGCM(AesCipher)
	if err != nil {
		loggers.Error("Error creating new GCM", "ERROR", err.Error())
		return nil, err
	}

	nonce := make([]byte, NewGcmBlock.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		loggers.Error("Error creating new nonce", "ERROR", err.Error())
		return nil, err
	}
	return NewGcmBlock.Seal(nonce, nonce, Data, nil), nil
}

type AesEncryption struct{}

func (a AesEncryption) Encrypter(data DomainLevel.IncomeEncryptData) ([]byte, error) {
	loggers := slog.With("EncryptAes")
	AesCipher, err := aes.NewCipher(data.Key)
	if err != nil {
		loggers.Error("Error creating new AesCipher", "ERROR", err.Error())
		return nil, err
	}
	NewGcmBlock, err := cipher.NewGCM(AesCipher)
	if err != nil {
		loggers.Error("Error creating new GCM", "ERROR", err.Error())
		return nil, err
	}

	nonce := make([]byte, NewGcmBlock.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		loggers.Error("Error creating new nonce", "ERROR", err.Error())
		return nil, err
	}
	return NewGcmBlock.Seal(nonce, nonce, data.Key, nil), nil
}

type RsaEncryption struct{}

func (r RsaEncryption) Encrypter(data DomainLevel.IncomeEncryptData) ([]byte, error) {
	keyData, err := x509.ParsePKCS1PublicKey(data.Key)
	if err != nil {
		slog.Error("Encrypter: error to parse a key", "ERROR", err)
		return nil, err
	}
	encryptAesKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, keyData, data.Data, nil)
	switch {
	case strings.Contains(fmt.Sprint(err), "message too long for RSA key size"):
		return nil, errors.New(DomainLevel.ErrorFileNameEncrypt)

	case err != nil:
		slog.Error("EncryptFileInfo; error to encrypt", "ERROR", err)
		return nil, errors.New(DomainLevel.ErrorStrangeCrypto)
	}
	return encryptAesKey, nil
}
