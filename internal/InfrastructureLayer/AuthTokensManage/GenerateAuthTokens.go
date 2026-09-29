package AuthTokensManage

import (
	"github.com/golang-jwt/jwt/v5"
)

type Generator interface {
	GenerateJWT(jwt.Claims) (string, error)
	GenerateRT(jwt.Claims) (string, error)
}

type CreatingTokens struct {
	Key []byte
}

func GetNewCreatingTokens() CreatingTokens {
	return CreatingTokens{}
}

func (c CreatingTokens) GenerateRT(claims jwt.Claims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return refreshToken.SignedString((c.Key))
}
func (c CreatingTokens) GenerateJWT(claims jwt.Claims) (string, error) {
	JwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return JwtToken.SignedString(c.Key)
}
