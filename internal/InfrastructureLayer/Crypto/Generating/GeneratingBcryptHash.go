package Generating

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

func (g Generating) GenerateHash(password []byte) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error generating hash from password", "Error", err.Error())
		return nil, err
	}
	return bytes, nil

}
