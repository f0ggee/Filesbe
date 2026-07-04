package Generating

import (
	"Kaban/internal/DomainLevel"
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

func (g Generating) GenerateHash(password []byte) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		slog.Error("GenerateHash; error to generate a hash", "ERROR", err)
		return nil, errors.New(DomainLevel.ErrorMakeHash)
	}
	return bytes, nil

}
