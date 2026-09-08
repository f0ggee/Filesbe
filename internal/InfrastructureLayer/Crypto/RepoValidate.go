package Crypto

import (
	"Kaban/internal/DomainLevel"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

const (
	ErrorSignCheck = "the signature isn't correct"
	ErrorPassword  = "the password isn't correct"
)

type Checking struct{}

func GetNeValidating() Checking {
	return Checking{}
}

func (c *Checking) PasswordVerify(hashOfPassword []byte, password []byte) error {
	err := bcrypt.CompareHashAndPassword((hashOfPassword), []byte(password))
	if err != nil {
		slog.Error("PasswordVerify;Error while checking the password", "Error", err.Error())
		return errors.New(ErrorPassword)

	}
	return nil
}

func (c *Checking) CheckSign(data DomainLevel.CheckSignKeyIncomingData) error {
	publicKeyMasterServer, err := x509.ParsePKCS1PublicKey(data.MasterPublicKey)
	if err != nil {
		slog.Error("CheckSignKey; Error marshalling public key", "ERROR", err.Error())
		return err
	}

	err = rsa.VerifyPKCS1v15(publicKeyMasterServer, crypto.SHA256, data.Hash, data.Sign)
	if err != nil {
		slog.Error("CheckSignKey; Error verifying signature", "Error", err.Error())
		return errors.New(ErrorSignCheck)
	}
	return nil
}
