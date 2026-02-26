package TokenInteraction

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func (A *ControlTokens) GenerateRT(claims jwt.Claims) (string, error) {

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(os.Getenv("KEYFORJWT")))
}
