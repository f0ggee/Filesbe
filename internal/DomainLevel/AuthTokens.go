package DomainLevel

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserAuthCheck struct {
	Jwt string
	Rft string
}
type ManageTokens interface {
	DeleteRefreshToken(string)
	SaveToken(string)
}

type Generator interface {
	GenerateJWT(jwt.Claims) (string, error)
	GenerateRT(jwt.Claims) (string, error)
}

type CheckingAuthTokens interface {
	CheckJwt(string) error
	CheckRt(string) (jwt.Claims, error)
	CheckingDenyList(string) bool
}

type AuthCheck interface {
	CheckUserAuth(UserAuthCheck) (string, error)
}
