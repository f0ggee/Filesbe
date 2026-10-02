package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"Kaban/internal/InfrastructureLayer/Tokens"
	"context"
	"crypto/rand"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RegisterApplication interface {
	RegisterApp(context.Context, *Dto.UserDataRegister) DomainLevel.RegisterApplicationOutComingData
}

type RegisterTest struct{}

func GetNewRegisterTest() *RegisterTest {
	return &RegisterTest{}
}

func (r RegisterTest) RegisterApp(ctx context.Context, register *Dto.UserDataRegister) DomainLevel.RegisterApplicationOutComingData {
	JwtToken := Tokens.GetNewJwtToken()
	RfToken := Tokens.GetNewRfToken()

	jwtMaker, err := JwtToken.Make()
	if err != nil {
		panic(err)
	}

	var bytes = []byte(rand.Text())
	jwtToken, err := jwtMaker.GetAuthToken(bytes)
	if err != nil {
		panic(err)
	}

	rfMaker, err := RfToken.Make()
	if err != nil {
		panic(err)
	}

	rfToken, err := rfMaker.GetAuthToken(bytes)
	if err != nil {
		panic(err)
	}

	return DomainLevel.RegisterApplicationOutComingData{
		Jwt: string(jwtToken),
		Rft: string(rfToken),
		Err: nil,
	}
}

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

func (sa *NewRegisterApplication) RegisterApp(ctx context.Context, de *Dto.UserDataRegister) DomainLevel.RegisterApplicationOutComingData {
	err := sa.CheckingDb.CheckerUser(ctx, de.Email)
	if err != nil {
		return DomainLevel.RegisterApplicationOutComingData{Err: err}
	}
	HashPassword, err := sa.Generator.GenerateHash([]byte(de.Password))
	if err != nil {
		return DomainLevel.RegisterApplicationOutComingData{Err: err}
	}

	UnitIdUser, err := sa.WriterDb.CreateUser(DomainLevel.CreateUserIncomingData{
		Name:         de.Name,
		Email:        de.Email,
		HashPassword: string(HashPassword),
		Ctx:          ctx,
	})
	if err != nil {
		return DomainLevel.RegisterApplicationOutComingData{
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
		return DomainLevel.RegisterApplicationOutComingData{Err: err}
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
		return DomainLevel.RegisterApplicationOutComingData{Err: err}
	}

	return DomainLevel.RegisterApplicationOutComingData{Rft: RefreshToken, Jwt: JwtToken}
}
