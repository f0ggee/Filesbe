package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"context"
	"crypto/rand"
	"errors"
	"log/slog"
	"time"

	"Kaban/internal/Dto"

	"github.com/golang-jwt/jwt/v5"
)

type NewLoginCrypto struct {
	Validate DomainLevel.CryptoValidating
}
type NewLoginData struct {
	ReaderDatabase DomainLevel.ReadDb
}
type NewLoginAuth struct {
	GeneratingTokens AuthTokensManage.Generator
}
type NewLogin struct {
	NewLoginData
	NewLoginCrypto
	NewLoginAuth
}

func GetNewLogin(newLoginData NewLoginData, newLoginCrypto NewLoginCrypto, newLoginAuth NewLoginAuth) *NewLogin {
	return &NewLogin{NewLoginData: newLoginData, NewLoginCrypto: newLoginCrypto, NewLoginAuth: newLoginAuth}
}

func (sa *NewLogin) LoginService(ctx context.Context, s Dto.UserLoginData) DomainLevel.LoginApplicationOutComingData {
	usersData := sa.ReaderDatabase.LoginData(s.Email, ctx)
	if usersData.Err != nil {
		return DomainLevel.LoginApplicationOutComingData{
			Err: usersData.Err,
		}
	}
	err := sa.Validate.PasswordVerify([]byte(usersData.HashPassword), []byte(s.Password))
	if err != nil {
		return DomainLevel.LoginApplicationOutComingData{
			Err: err,
		}
	}
	RefreshToken, err := sa.GeneratingTokens.GenerateRT(Dto.JwtCustomStruct{
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
		return DomainLevel.LoginApplicationOutComingData{
			Err: errors.New(DomainLevel.SessionError),
		}
	}
	JwtToken, err := sa.GeneratingTokens.GenerateJWT(Dto.JwtCustomStruct{
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
		return DomainLevel.LoginApplicationOutComingData{
			Err: errors.New(DomainLevel.SessionError),
		}
	}
	return DomainLevel.LoginApplicationOutComingData{
		Jwt: JwtToken,
		Rft: RefreshToken,
	}
}
