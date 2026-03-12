package DomainLevel

import (
	"Kaban/iternal/Dto"

	"github.com/golang-jwt/jwt/v5"
)

type ManageAuthTokens interface {
	GenerateJWT(jwt.Claims) (string, error)
	GenerateRT(jwt.Claims) (string, error)
	TokenDenyMapChecker(string) bool
	CheckLifeJwt(string) (*jwt.Token, error)
	CheckLifeRt(string) (*jwt.Token, error)
	DeleteAndSaveToken(string, string)
}

type ManageAuthTokensImpl interface {
	CollectDataForTokens(int) Dto.JwtCustomStruct
}
