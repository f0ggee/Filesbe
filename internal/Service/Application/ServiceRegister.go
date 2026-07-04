package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"context"
	"crypto/rand"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type NewRegisterApplication struct {
	DatabaseControlling
	Crypto
	AuthTokens
}

func GetNewNewRegisterApplication(databaseControlling DatabaseControlling, crypto Crypto, authTokens AuthTokens) *NewRegisterApplication {
	return &NewRegisterApplication{DatabaseControlling: databaseControlling, Crypto: crypto, AuthTokens: authTokens}
}

type RegisterApplicationOutComingData struct {
	Jwt string
	Rft string
	Err error
}

func (sa *NewRegisterApplication) RegisterService(de *Dto.UserDataRegister, ctx context.Context) RegisterApplicationOutComingData {

	err := sa.DatabaseControlling.Checker.CheckerUser(de.Email, ctx)
	if err != nil {
		return RegisterApplicationOutComingData{Err: err}
	}
	HashPassword, err := sa.Crypto.Generate.GenerateHash([]byte(de.Password))
	if err != nil {
		return RegisterApplicationOutComingData{Err: err}
	}

	UnitIdUser, err := sa.DatabaseControlling.Writer.CreateUser(DomainLevel.CreateUserIncomingData{
		Name:         de.Name,
		Email:        de.Email,
		HashPassword: string(HashPassword),
		Ctx:          ctx,
	})
	if err != nil {
		return RegisterApplicationOutComingData{
			Err: err,
		}
	}
	RefreshToken, err := sa.AuthTokens.GeneratingToken.GenerateRT(Dto.JwtCustomStruct{
		UserID: UnitIdUser,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Kabaner",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Hour)),
			ID:        rand.Text(),
		},
	})
	if err != nil {
		slog.Error("RegisterFunc; a strange error happened during creating a JWT token", "ERROR", err)
		return RegisterApplicationOutComingData{Err: err}
	}
	JwtToken, err := sa.AuthTokens.GeneratingToken.GenerateJWT(Dto.JwtCustomStruct{
		UserID: UnitIdUser,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Kabaner",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Hour)),
			ID:        rand.Text(),
		},
	})
	if err != nil {
		slog.Error("RegisterFunc; a strange error happened during creating a RFT token", "ERROR", err)
		return RegisterApplicationOutComingData{Err: err}
	}

	return RegisterApplicationOutComingData{Rft: RefreshToken, Jwt: JwtToken}
}
