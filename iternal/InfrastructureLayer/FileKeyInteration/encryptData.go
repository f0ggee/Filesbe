package FileKeyInteration

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

func (*FileInfoController) EncryptData(FileInfoData []byte, Key []byte) ([]byte, error) {
	keyRsa, err := x509.ParsePKCS1PrivateKey(Key)
	if err != nil {
		slog.Error("Func EncryptData ParsePKCS1PrivateKey fail", err)
		return nil, err
	}

	encryptAesKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &keyRsa.PublicKey, FileInfoData, nil)

	switch {
	case strings.Contains(fmt.Sprint(err), "message too long for RSA key size"):
		return nil, errors.New("file name is too long ")

	case err != nil:
		slog.Error("Error create aes", "ERR", err)
		return nil, errors.New("can't validate data ")

	}

	return encryptAesKey, nil
}
