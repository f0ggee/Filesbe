package TokenInteraction

import (
	"Kaban/iternal/Dto"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (S *ControlTokens) GenerateNewTokens(RT string) (string, string, error) {
	slog.Info("GenerateNewTokens", "Start", time.Now())

	key := []byte(os.Getenv("KEYFORJWT"))
	rt, err := jwt.ParseWithClaims(
		RT,
		&Dto.JwtCustomStruct{},
		func(token *jwt.Token) (any, error) {
			return key, nil
		},
	)

	if rt == nil {
		slog.Error("GenerateNewTokens", "Token is nil", err)
		return "", "", err
	}
	claims, ok := rt.Claims.(*Dto.JwtCustomStruct)
	if !ok {
		return "", "", errors.New("can't par")
	}

	slog.Info("GenerateNewTokens", "Claims", claims.RegisteredClaims.ID)
	jwtToken, err := S.GenerateJWT(claims.UserID)
	if err != nil {
		slog.Error("Err in CheckRt", err)
		return "", "", err
	}

	slog.Info("HelpfulType", "Token", claims.ExpiresAt)
	RtToken, err := S.GenerateRT(claims.UserID, claims.ExpiresAt)
	if err != nil {
		slog.Error("Error in checkRt", err)
		return "", "", err
	}
	slog.Info("GenerateNewTokens", "End", time.Now())

	return jwtToken, RtToken, nil
}
