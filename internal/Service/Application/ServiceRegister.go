package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
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

type NewRegisterApplication struct {
	CheckingDb DomainLevel.CheckingDb
	WriterDb   DomainLevel.WriteDb
	Generator  DomainLevel.CryptoGenerating

	Rft DomainLevel.AuthMaker
	Jwt DomainLevel.AuthMaker
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
	var BytesID [4]byte
	binary.BigEndian.PutUint32(BytesID[:], uint32(UnitIdUser))

	rfMaker, err := sa.Rft.Make()
	if err != nil {
		return DomainLevel.RegisterApplicationOutComingData{
			Jwt: "",
			Rft: "",
			Err: err,
		}
	}
	jwtMaker, err := sa.Jwt.Make()
	if err != nil {
		return DomainLevel.RegisterApplicationOutComingData{
			Jwt: "",
			Rft: "",
			Err: err,
		}
	}
	jwt, err := jwtMaker.GetAuthToken(BytesID[:])
	if err != nil {
		return DomainLevel.RegisterApplicationOutComingData{
			Jwt: "",
			Rft: "",
			Err: err,
		}
	}
	rfToken, err := rfMaker.GetAuthToken(BytesID[:])
	if err != nil {
		return DomainLevel.RegisterApplicationOutComingData{
			Jwt: "",
			Rft: "",
			Err: err,
		}
	}
	return DomainLevel.RegisterApplicationOutComingData{Rft: string(rfToken), Jwt: string(jwt)}
}
