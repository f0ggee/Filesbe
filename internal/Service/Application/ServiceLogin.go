package Application

import (
	"Kaban/internal/DomainLevel"
	"context"
	"crypto/rand"
	"errors"
	"log/slog"
	"time"

	"Kaban/internal/Dto"

	"github.com/golang-jwt/jwt/v5"
)

type NewLogin struct {
	DatabaseControlling
	GetCrypto
	AuthTokens
}
type LoginApplicationOutComingData struct {
	Jwt string
	Rft string
	Err error
}

func GetNewNewLogin(databaseControlling DatabaseControlling, crypto GetCrypto, authTokens AuthTokens) *NewLogin {
	return &NewLogin{DatabaseControlling: databaseControlling, GetCrypto: crypto, AuthTokens: authTokens}
}

func (sa *NewLogin) LoginService(s Dto.UserLoginData, ctx context.Context) LoginApplicationOutComingData {
	usersData := sa.DatabaseControlling.Reader.LoginData(s.Email, ctx)
	if usersData.Err != nil {
		return LoginApplicationOutComingData{
			Err: usersData.Err,
		}
	}
	err := sa.GetCrypto.Validate.PasswordVerify([]byte(usersData.HashPassword), []byte(s.Password))
	if err != nil {
		return LoginApplicationOutComingData{
			Err: err,
		}
	}
	RefreshToken, err := sa.AuthTokens.GeneratingToken.GenerateRT(Dto.JwtCustomStruct{
		UserID: (usersData.Id),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Kabaner",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Hour)),
			ID:        rand.Text(),
		},
	})
	if err != nil {
		slog.Error("LoginService; error to generate the Refresh Token", "ERROR", err)
		return LoginApplicationOutComingData{
			Err: errors.New(DomainLevel.ErrorCreateSession),
		}
	}
	JwtToken, err := sa.AuthTokens.GeneratingToken.GenerateJWT(Dto.JwtCustomStruct{
		UserID: usersData.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Kabaner",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Hour)),
			ID:        rand.Text(),
		},
	})
	if err != nil {
		slog.Error("LoginService; error to generate a Jwt token", "ERROR", err)
		return LoginApplicationOutComingData{
			Err: errors.New(DomainLevel.ErrorCreateSession),
		}
	}
	return LoginApplicationOutComingData{
		Jwt: JwtToken,
		Rft: RefreshToken,
	}
}
