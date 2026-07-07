package AuthTokensManage

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
)

type CheckingAuthTokens interface {
	CheckJwt(string) error
	CheckRt(string) (jwt.Claims, error)
}
type AuthCheck interface {
	CheckUserAuth(UserAuthCheckIncomingData) OutComingAuthData
}
type UserAuthCheckIncomingData struct {
	Jwt string
	Rft string
}
type OutComingAuthData struct {
	NewJwt          string
	IsNewJwtCreated bool
	Err             error
}
type NewAuthChecker struct {
	CreateTokens CreatingTokens
	Key          []byte
}

func GetNewNewAuthChecker(createTokens CreatingTokens, key []byte) *NewAuthChecker {
	return &NewAuthChecker{CreateTokens: createTokens, Key: key}
}

func (c NewAuthChecker) CheckRt(Rt string) (jwt.Claims, error) {
	Key, err := jwt.ParseWithClaims(Rt, &Dto.JwtCustomStruct{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return c.Key, nil
	})
	if err != nil {
		slog.Error("Error in parse refresh token", "error", err.Error())
		return nil, errors.New(DomainLevel.ErrorUserToken)
	}
	if !Key.Valid {
		return nil, errors.New(DomainLevel.ErrorUserNotAuthed)
	}
	return Key.Claims, nil
}

func (c NewAuthChecker) CheckJwt(JWT string) error {
	JwtToken, err := jwt.ParseWithClaims(JWT, &Dto.JwtCustomStruct{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return c.Key, nil
	})
	if err != nil {
		slog.Error("CheckJwt; error to parse a token", "ERROR", err.Error())
		return errors.New(DomainLevel.ErrorUserToken)
	}
	if !JwtToken.Valid {
		return errors.New(DomainLevel.ErrorUserToken)
	}
	return nil
}

func (n NewAuthChecker) CheckUserAuth(data UserAuthCheckIncomingData) OutComingAuthData {
	err := n.CheckJwt(data.Jwt)
	if err == nil {
		return OutComingAuthData{
			IsNewJwtCreated: true,
		}
	}
	Claims, err := n.CheckRt(data.Rft)
	if err != nil {
		return OutComingAuthData{
			Err: err,
		}
	}
	JwtToken, err := n.CreateTokens.GenerateJWT(Claims)
	if err != nil {
		return OutComingAuthData{Err: err}
	}
	return OutComingAuthData{
		NewJwt: JwtToken,
	}
}
