package manageTokensImpl

import (
	"Kaban/iternal/Dto"
	"crypto/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ManageTokensImpl struct{}

func (m ManageTokensImpl) CollectDataForTokens(id int) Dto.JwtCustomStruct {

	ds := Dto.JwtCustomStruct{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Kabaner",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			ID:        rand.Text(),
		},
	}

	return ds
}
