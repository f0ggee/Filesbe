package TokenInteraction

import (
	"Kaban/iternal/Dto"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (A *ControlTokens) GenerateRT(UnicId int, timeOldKey *jwt.NumericDate) (string, error) {
	timeLiveKey := &jwt.NumericDate{Time: time.Now().Add(24 * time.Hour)}

	if timeOldKey != nil {
		timeLiveKey = timeOldKey

	}

	claims := Dto.JwtCustomStruct{
		UserID: UnicId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: timeLiveKey,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Admin",
		},
	}
	token1 := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	sd := []byte(os.Getenv("KEYFORJWT"))
	tokenString, err := token1.SignedString(sd)
	if err != nil {
		slog.Error("Error sign cookie", err)
		return "", err
	}

	return tokenString, nil
}
