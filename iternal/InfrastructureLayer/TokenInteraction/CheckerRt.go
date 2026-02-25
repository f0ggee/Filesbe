package TokenInteraction

import (
	"Kaban/iternal/Dto"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (A *ControlTokens) CheckLifeRt(Rt string) error {

	key := []byte(os.Getenv("KEYFORJWT"))
	_, err := jwt.ParseWithClaims(Rt, &Dto.JwtCustomStruct{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return key, nil // Возвращаем секретный ключ для проверки
	})

	if err != nil {
		delete(Dto.AllowList, Rt)
		Dto.DenyList[Rt] = time.Now()
		slog.Error("Refresh have been change", err)

	}

	return nil

}
