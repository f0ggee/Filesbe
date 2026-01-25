package DomainLevel

import (
	"github.com/golang-jwt/jwt/v5"
)

type ManageTokens interface {
	GenerateJWT(int) (string, error)
	GenerateRT(int, *jwt.NumericDate) (string, error)
	TokenDenyMapChecker(string) bool
	CheckLifeJwt(string) error
	CheckLifeRt(string) error
	GenerateNewTokens(string) (string, string, error)
	DeleteAndSaveToken(string, string)
}
