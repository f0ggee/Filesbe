package Crypto

import (
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

type CheckingTest struct {
}

func NewCheckingTest() *CheckingTest {
	return &CheckingTest{}
}

func (c CheckingTest) PasswordVerify(bytes []byte, bytes2 []byte) error {
	return nil
}

type CheckingTestError struct {
}

func NewCheckingTestError() *CheckingTestError {
	return &CheckingTestError{}
}

func (c CheckingTestError) PasswordVerify(bytes []byte, bytes2 []byte) error {

	return ErrorPassword
}
