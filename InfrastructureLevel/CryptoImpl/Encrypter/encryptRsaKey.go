package Encrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"log/slog"
)

func (s *Encryption) EncryptRsaKey(AesKey []byte, RsaKey []byte) ([]byte, error) {
	aesNewBlock, err := aes.NewCipher(AesKey)
	if err != nil {
		slog.Error("error creating new aes cipher", "Error", err.Error())
		return nil, err

	}
	gcm, err := cipher.NewGCM(aesNewBlock)
	if err != nil {
		slog.Error("error creating new gcm", "Error", err.Error())
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())

	text := gcm.Seal(nonce, nonce, RsaKey, nil)

	return text, nil
}
