package Crypto

import (
	"Kaban/internal/DomainLevel"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

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
