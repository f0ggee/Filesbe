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
	databaseControlling
	getCrypto
	authTokens
}

func GetNewNewRegisterApplication(databaseControlling databaseControlling, crypto getCrypto, authTokens authTokens) *NewRegisterApplication {
	return &NewRegisterApplication{databaseControlling: databaseControlling, getCrypto: crypto, authTokens: authTokens}
}

type RegisterApplicationOutComingData struct {
	Jwt string
	Rft string
	Err error
}

func (sa *NewRegisterApplication) RegisterService(de *Dto.UserDataRegister, ctx context.Context) RegisterApplicationOutComingData {

	err := sa.databaseControlling.Checker.CheckerUser(de.Email, ctx)
	if err != nil {
		return RegisterApplicationOutComingData{Err: err}
	}
	HashPassword, err := sa.getCrypto.Generate.GenerateHash([]byte(de.Password))
	if err != nil {
		return RegisterApplicationOutComingData{Err: err}
	}

	UnitIdUser, err := sa.databaseControlling.Writer.CreateUser(DomainLevel.CreateUserIncomingData{
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
	RefreshToken, err := sa.authTokens.GeneratingToken.GenerateRT(Dto.JwtCustomStruct{
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
	JwtToken, err := sa.authTokens.GeneratingToken.GenerateJWT(Dto.JwtCustomStruct{
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
