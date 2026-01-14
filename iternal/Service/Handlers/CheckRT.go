package Handlers

import (
	"Kaban/iternal/Dto"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func CheckRefreshTokenLifeTime(RT string) error {

	key := []byte(os.Getenv("KEYFORJWT"))
	_, err := jwt.ParseWithClaims(RT, &Dto.MyCustomCookie{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка алгоритма (важно для безопасности, чтобы избежать атак alg:none)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return key, nil // Возвращаем секретный ключ для проверки
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return errors.New("refresh Token expired")
		}

		slog.Error("Error check  RT token", err)
		slog.Info("RT", RT)

		return err
	}
	return nil

}

func GenerateNewTokens(RT string) (string, string, int, string, error) {

	key := []byte(os.Getenv("KEYFORJWT"))
	token, err := jwt.ParseWithClaims(
		RT,
		&Dto.MyCustomCookie{},
		func(token *jwt.Token) (any, error) {
			return key, nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", "", 0, "", errors.New("refresh Token expired")
		}

		slog.Error("Error check  RT token", err)

		return "", "", 0, "", err
	}

	claims, ok := token.Claims.(*Dto.MyCustomCookie)
	if !ok {
		return "", "", 0, "", errors.New("can't par")
	}

	jwtToken, err := JWT(claims.UserID, claims.Sc)
	if err != nil {
		slog.Error("Err in CheckRt", err)
		return "", "", 0, "", err
	}

	RtToken, err := RFT(claims.UserID, claims.Sc)
	if err != nil {
		slog.Error("Erro in checkRt", err)
		return "", "", 0, "", err
	}

	return jwtToken, RtToken, claims.UserID, claims.Sc, err
}
