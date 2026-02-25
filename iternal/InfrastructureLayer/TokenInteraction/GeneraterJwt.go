package TokenInteraction

import (
	Dto2 "Kaban/iternal/Dto"
	//Dto2 "Kaban/iternal/Dto"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (A *ControlTokens) GenerateJWT(id int) (string, error) {
	claims := Dto2.JwtCustomStruct{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Admin",
		},
	}
	token1 := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	sd := []byte(os.Getenv("KEYFORJWT"))

	tokenString, err := token1.SignedString(sd)
	if err != nil {
		slog.Error("Error sign cookie jwt", err)
		return "", err
	}

	return tokenString, nil

}
