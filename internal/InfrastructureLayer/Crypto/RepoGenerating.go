package Crypto

import (
	"crypto/rand"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type Generating struct{}

func GetNewGenerating() Generating {
	return Generating{}
}

func (g Generating) GenerateText(n int) string {
	if n == 0 && n > len(rand.Text()) {
		return rand.Text()
	}
	return rand.Text()[:5]
}
func (g Generating) GenerateHash(password []byte) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		slog.Error("GenerateHash; error to generate a hash", "ERROR", err)
		return nil, ErrorMakeHash
	}
	return bytes, nil

}
