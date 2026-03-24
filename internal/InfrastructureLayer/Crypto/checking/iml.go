package checking

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type Checking struct{}

func (c *Checking) PasswordVerify(hashOfPassword []byte, password []byte) error {
	slog.Info("Password checking starts")
	err := bcrypt.CompareHashAndPassword([]byte(hashOfPassword), []byte(password))
	if err != nil {
		slog.Error("Error while checking the password", "Error", err.Error())
		return err

	}
	slog.Info("Password ends")
	return nil
}
