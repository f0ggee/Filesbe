package Crypto

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/awnumar/memguard"
)

type Decryption struct {
	RepoParsers.Decode
}

func GetNewDecryption(decode RepoParsers.Decode) *Decryption {
	return &Decryption{Decode: decode}
}

func (d Decryption) SayHello(string) string {

	return "Hello World"
}
func (d Decryption) DecryptPacket(aesKey []byte, plainText []byte) *memguard.LockedBuffer {
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		slog.Error("Func DecryptPacket:Error create new aes block", "Error", err.Error())
		return nil
	}
	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		slog.Error("Func DecryptPacket: Error create new gcm", "Error", err.Error())
		return nil
	}
	packetData, err := gcm.Open(nil, plainText[:gcm.NonceSize()], plainText[gcm.NonceSize():], nil)

	if err != nil {
		slog.Error("Func DecryptPacket: Error decrypt packet", "Error", err.Error())
		return nil
	}
	lockedBuffer := memguard.NewBufferFromBytes(packetData)
	defer memguard.WipeBytes(packetData)
	return lockedBuffer
}

func (d Decryption) DecryptAesKey(RsaKey []byte, aesKey []byte) ([]byte, error) {
	RsaKeyPrivate, err := x509.ParsePKCS1PrivateKey(RsaKey)
	if err != nil {
		slog.Error("DecryptAesKey;Error Parsing RsaKey", "Func decrypt error", err)
		return nil, errors.New(DomainLevel.ErrorDecryptKeys)
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, RsaKeyPrivate, aesKey, nil)
}

func (d Decryption) DecryptFileInfo(FileInfo []byte, NewRsaKey []byte, OldRsaKey []byte) ([]byte, string, error) {
	keyRsa, err := x509.ParsePKCS1PrivateKey(NewRsaKey)
	if err != nil {
		slog.Error("Func DecryptFileInfo ParsePKCS1PrivateKey fail", "ERROR", err)
		return nil, "", errors.New(DomainLevel.ErrorDecryptFileInfo)
	}
	decryptFileInfo, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, keyRsa, FileInfo, nil)
	switch {
	case strings.Contains(fmt.Sprint(err), "decryption error"):
		keyRsaOld, err := x509.ParsePKCS1PrivateKey(OldRsaKey)
		if err != nil {
			slog.Error("Func DecryptFileInfo ParsePKCS1PrivateKey fail", "ERROR", err.Error())
			return nil, "", errors.New(DomainLevel.ErrorInfoOld)
		}
		decryptFileInfo, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, keyRsaOld, FileInfo, nil)
		if err != nil {
			slog.Error("DecryptFileIno;error decrypt file info", "ERROR", err.Error())
			return nil, "", errors.New(DomainLevel.ErrorDecryptFileInfo)
		}

	}
	sa := &DomainLevel.FileLabelsBytes{
		FileName: "",
		AesKey:   "",
	}

	err = d.Decode.JsonDecodeMarshall(&sa, decryptFileInfo)
	if err != nil {
		return nil, "", err
	}
	aesKeyIntoByte, err := hex.DecodeString(sa.AesKey)
	if err != nil {
		slog.Error("Func DecryptFileInfo;error decode aes key into string", "Error", err.Error())
		return nil, "", errors.New(DomainLevel.ErrorGetFileInfo)
	}

	return aesKeyIntoByte, sa.FileName, nil
}
