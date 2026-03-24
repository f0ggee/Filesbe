package Handlers

import (
	"Kaban/internal/Dto"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (sa *HandlerPackCollect) RegisterService(de *Dto.UserDataRegister, ctx context.Context) (string, string, error) {

	err := sa.DatabaseControlling.Checker.CheckerUser(de.Email, ctx)
	switch {
	case errors.Is(err, errors.New("person already exist")):
		return "", "", errors.New("person already exist")

	case err != nil:
		return "", "", err
	}
	HashPassword, err := sa.Crypto.Generate.GenerateHashFromPassword([]byte(de.Password))
	if err != nil {
		slog.Error("Err generate a password-scrypt", "err", err)
		return "", "", err
	}

	UnitIdUser, err := sa.DatabaseControlling.Writer.CreateUser(de.Name, de.Email, hex.EncodeToString(HashPassword), ctx)
	if err != nil {
		return "", "", err
	}

	RefreshToken, err := sa.AuthTokens.GeneratingToken.GenerateRT(Dto.JwtCustomStruct{
		UserID: UnitIdUser,
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
		UserID: UnitIdUser,
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

	return JwtToken, RefreshToken, nil
}
