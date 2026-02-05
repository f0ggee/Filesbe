package FileKeyInteration

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"log/slog"
)

func (*FileInfoController) EncryptData(FileInfoData []byte, Key *rsa.PublicKey) ([]byte, error) {
	encryptAesKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, Key, FileInfoData, nil)
	if err != nil {
		slog.Error("Error create aes", "ERR", err)
		return nil, errors.New("can't validate data ")
	}

	return encryptAesKey, nil
}
