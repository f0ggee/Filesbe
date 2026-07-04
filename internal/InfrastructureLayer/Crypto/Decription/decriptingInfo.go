package Decription

import (
	"Kaban/internal/Dto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

const (
	ErrorDecryptFileInfo = "an error happened during decrypting info"
	ErrorInfoOld         = "the file info isn't valid"
	ErrorGetFileInfo     = "can't parse file info"
)

func (d DecryptionData) DecryptFileInfo(FileInfo []byte, NewRsaKey []byte, OldRsaKey []byte) ([]byte, string, error) {
	keyRsa, err := x509.ParsePKCS1PrivateKey(NewRsaKey)
	if err != nil {
		slog.Error("Func DecryptFileInfo ParsePKCS1PrivateKey fail", "Error", err)
		return nil, "", errors.New(ErrorDecryptFileInfo)
	}
	decryptFileInfo, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, keyRsa, FileInfo, nil)

	switch {
	case strings.Contains(fmt.Sprint(err), "decryption error"):
		keyRsaOld, err := x509.ParsePKCS1PrivateKey(OldRsaKey)
		if err != nil {
			slog.Error("Func DecryptFileInfo ParsePKCS1PrivateKey fail", "ERROR", err.Error())
			return nil, "", errors.New(ErrorInfoOld)
		}
		decryptFileInfo, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, keyRsaOld, FileInfo, nil)
		if err != nil {
			slog.Error("Error also decrypt with an old key", "ERROR", err.Error())
			return nil, "", errors.New(ErrorDecryptFileInfo)
		}

	}

	sa := &Dto.FileLabelsBytes{
		FileName: "",
		AesKey:   "",
	}
	err = json.Unmarshal(decryptFileInfo, &sa)
	if err != nil {
		slog.Error("Error unmarshal aes", "ERR", err)
		return nil, "", errors.New(ErrorGetFileInfo)
	}

	aesKeyIntoByte, err := hex.DecodeString(sa.AesKey)
	if err != nil {
		slog.Error("Func DecryptFileInfo;error decode aes key into string", "Error", err.Error())
		return nil, "", errors.New(ErrorGetFileInfo)
	}

	return aesKeyIntoByte, sa.FileName, nil
}
