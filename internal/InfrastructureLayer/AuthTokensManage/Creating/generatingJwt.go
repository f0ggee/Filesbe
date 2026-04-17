package Creating

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func (c CreatingTokens) GenerateJWT(claims jwt.Claims) (string, error) {
	JwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return JwtToken.SignedString([]byte(os.Getenv("KEYFORJWT")))

}
