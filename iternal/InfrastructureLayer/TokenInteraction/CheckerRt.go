package TokenInteraction

import (
	"Kaban/iternal/Dto"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func (A *ControlTokens) CheckLifeRt(Rt string) (*jwt.Token, error) {

	key := []byte(os.Getenv("KEYFORJWT"))
	Key, err := jwt.ParseWithClaims(Rt, &Dto.JwtCustomStruct{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return key, nil // Возвращаем секретный ключ для проверки
	})

	if err != nil {
		slog.Error("Error in parse refresh token", "error", err.Error())
		return nil, err
	}

	return Key, nil

}
