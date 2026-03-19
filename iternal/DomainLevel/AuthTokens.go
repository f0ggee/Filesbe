package DomainLevel

import (
	"github.com/golang-jwt/jwt/v5"
)

type ManageTokens interface {
	DeleteRefreshToken(string)
	SaveToken(string)
}

type Generator interface {
	GenerateJWT(jwt.Claims) (string, error)
	GenerateRT(jwt.Claims) (string, error)
}

type CheckingAuthTokens interface {
	CheckJwt(string) (*jwt.Token, error)
	CheckRt(string) (*jwt.Token, error)
	CheckingDenyList(string) bool
}
