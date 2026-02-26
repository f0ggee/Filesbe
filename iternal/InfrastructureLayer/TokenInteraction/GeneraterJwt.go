package TokenInteraction

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func (A *ControlTokens) GenerateJWT(claims jwt.Claims) (string, error) {

	JwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return JwtToken.SignedString([]byte(os.Getenv("KEYFORJWT")))

}
