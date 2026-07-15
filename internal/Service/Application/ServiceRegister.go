package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"context"
	"crypto/rand"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type NewRegisterCrypto struct {
	Generator DomainLevel.CryptoGenerating
}
type NewRegisterDataMange struct {
	CheckingDb      DomainLevel.CheckingDb
	WriterDb        DomainLevel.WriteDb
	GeneratorTokens AuthTokensManage.Generator
}
type NewRegisterApplication struct {
	NewRegisterDataMange
	NewRegisterCrypto
}

func GetNewRegisterApplication(newRegisterDataMange NewRegisterDataMange, newRegisterCrypto NewRegisterCrypto) *NewRegisterApplication {
	return &NewRegisterApplication{NewRegisterDataMange: newRegisterDataMange, NewRegisterCrypto: newRegisterCrypto}
}

type RegisterApplicationOutComingData struct {
	Jwt string
	Rft string
	Err error
}

func (sa *NewRegisterApplication) RegisterService(de *Dto.UserDataRegister, ctx context.Context) RegisterApplicationOutComingData {

	err := sa.CheckingDb.CheckerUser(de.Email, ctx)
	if err != nil {
		return RegisterApplicationOutComingData{Err: err}
	}
	HashPassword, err := sa.Generator.GenerateHash([]byte(de.Password))
	if err != nil {
		return RegisterApplicationOutComingData{Err: err}
	}

	UnitIdUser, err := sa.WriterDb.CreateUser(DomainLevel.CreateUserIncomingData{
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
	RefreshToken, err := sa.GeneratorTokens.GenerateRT(Dto.JwtCustomStruct{
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
	JwtToken, err := sa.GeneratorTokens.GenerateJWT(Dto.JwtCustomStruct{
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
