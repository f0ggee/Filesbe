package Crypto

import (
	"Kaban/internal/DomainLevel"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"log/slog"
)

type RsaSigner struct {
	key []byte
}

func NewRsaSigner() *RsaSigner {
	return &RsaSigner{}
}

func (r *RsaSigner) MakeSign(signData []byte) ([]byte, error) {
	KeyPrivate, err := x509.ParsePKCS1PrivateKey(r.key)
	if err != nil {
		slog.Error("RsaMake: error to convert a key", "ERROR", err)
		return nil, ErrorStrangeCrypto
	}
	HashData := sha256.Sum256(signData)
	SignedMessage, err := rsa.SignPKCS1v15(rand.Reader, KeyPrivate, crypto.SHA256, HashData[:])
	if err != nil {
		slog.Error("GenerateSignature;Error while signing message", "ERROR", err)
		return nil, ErrorMakeSign
	}
	return SignedMessage, nil
}

func (r *RsaSigner) CheckSign(signature []byte) error {
	publicKeyMasterServer, err := x509.ParsePKCS1PublicKey(r.key)
	if err != nil {
		slog.Error("CheckSignKey; Error marshalling public key1", "ERROR", err.Error())
		return err
	}
	hash := sha256.New().Sum(signature)
	err = rsa.VerifyPKCS1v15(publicKeyMasterServer, crypto.SHA256, hash, signature)
	if err != nil {
		slog.Error("CheckSignKey; Error verifying signature", "Error", err.Error())
		return ErrorSignCheck
	}

	return nil
}

func (r *RsaSigner) SetAdditionalData(data []byte) DomainLevel.SignatureMaker {
	return r
}

func (r *RsaSigner) MakeSigner(key []byte) (DomainLevel.Signer, error) {
	r.key = key
	return r, nil
}
