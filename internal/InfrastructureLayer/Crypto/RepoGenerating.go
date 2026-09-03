package Crypto

import (
	"Kaban/internal/DomainLevel"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type Generating struct{}

func GetNewGenerating() Generating {
	return Generating{}
}

func (g Generating) GenerateShortName() string {
	NewString := rand.Text()
	return NewString[:5]
}
func (g Generating) GenerateSignature(message []byte, key []byte) ([]byte, error) {
	KeyPrivate, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		slog.Error("GenerateSignature;Error while converting key to private key", "ERROR", err)
		return nil, errors.New(DomainLevel.ErrorStrangeCrypto)
	}
	HashData := sha256.Sum256(message)
	SignedMessage, err := rsa.SignPKCS1v15(rand.Reader, KeyPrivate, crypto.SHA256, HashData[:])
	if err != nil {
		slog.Error("GenerateSignature;Error while signing message", "ERROR", err)
		return nil, errors.New(DomainLevel.ErrorMakeSign)
	}
	return SignedMessage, nil

}
func (g Generating) GenerateHash(password []byte) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		slog.Error("GenerateHash; error to generate a hash", "ERROR", err)
		return nil, errors.New(DomainLevel.ErrorMakeHash)
	}
	return bytes, nil

}
