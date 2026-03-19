package Encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log/slog"
)

func (e *Encrypter) EncryptAes(AesKey []byte, Data []byte) ([]byte, error) {
	AesCipher, err := aes.NewCipher(AesKey)
	if err != nil {
		slog.Error("Error creating new AesCipher", "Error", err.Error())
		return nil, err
	}

	NewGcmBlock, err := cipher.NewGCM(AesCipher)
	if err != nil {
		slog.Error("Error creating new GCM", "Error", err.Error())
		return nil, err
	}

	nonce := make([]byte, NewGcmBlock.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		slog.Error("Error creating new nonce", "Error", err.Error())
		return nil, err
	}

	return NewGcmBlock.Seal(nonce, nonce, Data, nil), nil

}
