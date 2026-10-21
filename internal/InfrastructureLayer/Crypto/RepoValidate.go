package Crypto

import (
	"Kaban/internal/DomainLevel"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type Checking struct{}

func GetNewChecking() *Checking {
	return &Checking{}
}

func (c *Checking) PasswordVerify(hashOfPassword []byte, password []byte) error {
	err := bcrypt.CompareHashAndPassword(hashOfPassword, password)
	if err != nil {
		slog.Error("PasswordVerify: the password isn't correct", "ERROR", err.Error())
		return ErrorPassword

	}
	return nil
}

func (c *Checking) CheckSign(data DomainLevel.CheckSignKeyIncomingData) error {
	publicKeyMasterServer, err := x509.ParsePKCS1PublicKey(data.MasterPublicKey)
	if err != nil {
		slog.Error("CheckSignKey; Error marshalling public key1", "ERROR", err.Error())
		return err
	}

	err = rsa.VerifyPKCS1v15(publicKeyMasterServer, crypto.SHA256, data.Hash, data.Sign)
	if err != nil {
		slog.Error("CheckSignKey; Error verifying signature", "Error", err.Error())
		return ErrorSignCheck
	}
	return nil
}

type CheckingTest struct {
}

func NewCheckingTest() *CheckingTest {
	return &CheckingTest{}
}

func (c CheckingTest) CheckSign(data DomainLevel.CheckSignKeyIncomingData) error {
	return nil
}

func (c CheckingTest) PasswordVerify(bytes []byte, bytes2 []byte) error {
	return nil
}

type CheckingTestError struct {
}

func NewCheckingTestError() *CheckingTestError {
	return &CheckingTestError{}
}

func (c CheckingTestError) CheckSign(data DomainLevel.CheckSignKeyIncomingData) error {
	return nil
}

func (c CheckingTestError) PasswordVerify(bytes []byte, bytes2 []byte) error {

	return ErrorPassword
}
