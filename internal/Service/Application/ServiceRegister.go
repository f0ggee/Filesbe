package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"Kaban/internal/InfrastructureLayer/Tokens"
	"context"
	"crypto/rand"
	"encoding/binary"

	"github.com/awnumar/memguard"
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
	TokenMaker1 DomainLevel.AuthMaker
	TokenMaker2 DomainLevel.AuthMaker
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
	memguard.WipeBytes([]byte(de.Password))
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

	jwtMaker, err := sa.TokenMaker1.Make()
	if err != nil {
		return DomainLevel.RegisterApplicationOutComingData{
			Jwt: "",
			Rft: "",
			Err: err,
		}
	}
	var BytesID [4]byte
	binary.BigEndian.PutUint32(BytesID[:], uint32(UnitIdUser))
	jwtToken, err := jwtMaker.GetAuthToken(BytesID[:])
	if err != nil {
		return DomainLevel.RegisterApplicationOutComingData{Err: err}
	}
	refreshMaker, err := sa.TokenMaker2.Make()
	if err != nil {
		return DomainLevel.RegisterApplicationOutComingData{
			Jwt: "",
			Rft: "",
			Err: err,
		}
	}

	refreshToken, err := refreshMaker.GetAuthToken(BytesID[:])
	return DomainLevel.RegisterApplicationOutComingData{Rft: string(refreshToken), Jwt: string(jwtToken)}
}
