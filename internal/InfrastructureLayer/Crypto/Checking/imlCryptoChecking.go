package Checking

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type Validating struct{}

func (c *Validating) PasswordVerify(hashOfPassword []byte, password []byte) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashOfPassword), []byte(password))
	if err != nil {
		slog.Error("Error while Validator the password", "Error", err.Error())
		return err

	}
	return nil
}
