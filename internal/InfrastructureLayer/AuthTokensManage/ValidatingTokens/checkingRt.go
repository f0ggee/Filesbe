package ValidatingTokens

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
)

func (c Checking) CheckRt(Rt string) (jwt.Claims, error) {
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
