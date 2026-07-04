package ValidatingTokens

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
)

func (c Checking) CheckJwt(JWT string) error {
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
