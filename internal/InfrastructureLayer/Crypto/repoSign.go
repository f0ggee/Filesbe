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

type RsaSign struct {
	key []byte
}

func NewRsaSigner() *RsaSign {
	return &RsaSign{}
}

func (r *RsaSign) MakeSign(signData []byte) ([]byte, error) {
	KeyPrivate, err := x509.ParsePKCS1PrivateKey(r.key)
	if err != nil {
		slog.Error("RsaMake: error to convert a key", "ERROR", err)
		return nil, ErrorStrangeCrypto
	}
	HashData := sha256.Sum256(signData)
	SignedMessage, err := rsa.SignPKCS1v15(rand.Reader, KeyPrivate, crypto.SHA256, HashData[:])
	if err != nil {
		slog.Error("GenerateSignature: error while signing message", "ERROR", err)
		return nil, ErrorMakeSign
	}
	return SignedMessage, nil
}

func (r *RsaSign) CheckSign(signature []byte, payload []byte) error {

	var masterKey *rsa.PublicKey

	if _masterKey, err := x509.ParsePKCS1PrivateKey(r.key); err == nil {
		masterKey = &_masterKey.PublicKey
	} else {
		_masterKeyP, _err := x509.ParsePKCS1PublicKey(r.key)
		if _err != nil {
			return ErrorInvalidData
		}
		masterKey = _masterKeyP
	}

	hash := sha256.Sum256(payload)
	err := rsa.VerifyPKCS1v15(masterKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		slog.Error("CheckSignKey; Error verifying signature", "Error", err.Error())
		return ErrorSignCheck
	}

	return nil
}

func (r *RsaSign) SetAdditionalData(data []byte) DomainLevel.SignatureMaker {
	return r
}

func (r *RsaSign) MakeSigner(key []byte) (DomainLevel.Signer, error) {
	r.key = key
	return r, nil
}
