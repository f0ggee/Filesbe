package Handlers

import (
	"Kaban/Dto"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"os"
)

func CheckLifeJWT(JWT string) (*jwt.Token, error) {

	key := []byte(os.Getenv("KEYFORJWT"))
	token, err := jwt.ParseWithClaims(JWT, &Dto.MyCustomCookie{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка алгоритма (важно для безопасности, чтобы избежать атак alg:none)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return key, nil // Возвращаем секретный ключ для проверки
	})

	if err != nil {
		slog.Error("ERRR", JWT)
		slog.Error("Error in check Jwt", err)
		return nil, err
	}
	return token, nil

}
