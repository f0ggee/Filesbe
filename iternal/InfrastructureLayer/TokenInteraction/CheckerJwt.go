package TokenInteraction

import (
	"Kaban/iternal/Dto"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func (A *ControlTokens) CheckLifeJwt(JWT string) error {
	key := []byte(os.Getenv("KEYFORJWT"))
	_, err := jwt.ParseWithClaims(JWT, &Dto.JwtCustomStruct{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		slog.Error("Error in check Jwt", err)
		return err
	}
	return nil

}
