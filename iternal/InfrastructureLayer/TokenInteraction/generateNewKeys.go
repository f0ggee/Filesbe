package TokenInteraction

import (
	"Kaban/iternal/Dto"
	"errors"
	"log/slog"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func (S *ControlTokens) GenerateNewTokens(RT string) (string, string, error) {

	key := []byte(os.Getenv("KEYFORJWT"))
	rt, err := jwt.ParseWithClaims(
		RT,
		&Dto.JwtCustomStruct{},
		func(token *jwt.Token) (any, error) {
			return key, nil
		},
	)

	claims, ok := rt.Claims.(*Dto.JwtCustomStruct)
	if !ok {
		return "", "", errors.New("can't par")
	}

	jwtToken, err := S.GenerateJWT(claims.UserID)
	if err != nil {
		slog.Error("Err in CheckRt", err)
		return "", "", err
	}

	RtToken, err := S.GenerateRT(claims.UserID, claims.ExpiresAt)
	if err != nil {
		slog.Error("Error in checkRt", err)
		return "", "", err
	}

	return jwtToken, RtToken, nil
}
