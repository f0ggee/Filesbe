package ValidatingTokens

import (
	"Kaban/internal/Dto"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func (c Checking) CheckJwt(JWT string) (*jwt.Token, error) {
	key := []byte(os.Getenv("KEYFORJWT"))
	JwtToken, err := jwt.ParseWithClaims(JWT, &Dto.JwtCustomStruct{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		slog.Error("Erorr to parse jwt token", "Error", err.Error())
		return nil, err
	}

	return JwtToken, nil
}
