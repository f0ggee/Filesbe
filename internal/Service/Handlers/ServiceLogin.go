package Handlers

import (
	Dto2 "Kaban/internal/Dto"
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"time"

	"Kaban/internal/Dto"

	"github.com/golang-jwt/jwt/v5"
)

func (sa *HandlerPackCollect) LoginService(s Dto2.UserLoginData, ctx context.Context) (string, string, error) {

	slog.Info("Func LoginService starts")

	Id, password, err := sa.DatabaseControlling.Reader.LoginData(s.Email, ctx)

	if err != nil {
		slog.Error("Error in LoginData", "error", err)
		return "", "", err
	}

	PasswordBytes, err := hex.DecodeString(password)
	if err != nil {
		slog.Error("func decoding login user's password", "err", err)
		return "", "", err
	}
	err = sa.Crypto.Validate.PasswordVerify([]byte(password), PasswordBytes)
	if err != nil {
		return "", "", err
	}

	RefreshToken, err := sa.AuthTokens.GeneratingToken.GenerateRT(Dto.JwtCustomStruct{
		UserID: Id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Kabaner",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Hour)),
			ID:        rand.Text(),
		},
	})
	if err != nil {
		slog.Error("func login 3", "err", err)
		return "", "", err
	}
	JwtToken, err := sa.AuthTokens.GeneratingToken.GenerateJWT(Dto.JwtCustomStruct{
		UserID: Id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Kabaner",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Hour)),
			ID:        rand.Text(),
		},
	})
	if err != nil {
		slog.Error("func login 4", "err", err)
		return "", "", err
	}

	slog.Info("Func LoginService ends")
	return JwtToken, RefreshToken, nil

}
