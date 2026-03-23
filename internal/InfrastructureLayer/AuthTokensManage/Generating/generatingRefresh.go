package Generating

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func (c CreatingTokens) GenerateRT(claims jwt.Claims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(os.Getenv("KEYFORJWT")))
}
