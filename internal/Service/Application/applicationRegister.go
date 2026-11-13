package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Drivers/Tokens"
	"crypto/rand"
	"encoding/binary"

	"github.com/awnumar/memguard"
)

type RegisterTest struct{}

func GetNewRegisterTest() *RegisterTest {
	return &RegisterTest{}
}

func (r RegisterTest) RegisterApplication(input InputRegister) DomainLevel.RegisterApplicationOutComingData {
	JwtToken := Tokens.NewJwtToken()
	RfToken := Tokens.NewRfToken()

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

func (sa *NewRegisterApplication) RegisterApplication(input InputRegister) DomainLevel.RegisterApplicationOutComingData {
	defer memguard.WipeBytes(input.Password)
	err := sa.CheckingDb.CheckerUser(input.Ctx, input.Email)
	if err != nil {
		return DomainLevel.RegisterApplicationOutComingData{Err: err}
	}
	HashPassword, err := sa.Generator.GenerateHash(input.Password)
	if err != nil {
		return DomainLevel.RegisterApplicationOutComingData{Err: err}
	}
	memguard.WipeBytes(input.Password)
	UnitIdUser, err := sa.WriterDb.CreateUser(DomainLevel.CreateUserIncomingData{
		Name:         input.Name,
		Email:        input.Email,
		HashPassword: string(HashPassword),
		Ctx:          input.Ctx,
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
