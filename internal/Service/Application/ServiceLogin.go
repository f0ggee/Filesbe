package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"context"
	"encoding/binary"
)

type LoginApplication interface {
	Login(context.Context, Dto.UserLoginData) DomainLevel.LoginApplicationOutComingData
}
type NewLoginCrypto struct {
	Validate DomainLevel.CryptoValidating
}
type NewLoginData struct {
	ReaderDatabase DomainLevel.ReadDb
}
type NewLoginAuth struct {
	GenereteTokens1 DomainLevel.AuthMaker
	GenereteToken2  DomainLevel.AuthMaker
}
type NewLogin struct {
	NewLoginData
	NewLoginCrypto
	NewLoginAuth
}

func GetNewLogin(newLoginData NewLoginData, newLoginCrypto NewLoginCrypto, newLoginAuth NewLoginAuth) *NewLogin {
	return &NewLogin{NewLoginData: newLoginData, NewLoginCrypto: newLoginCrypto, NewLoginAuth: newLoginAuth}
}

func (sa *NewLogin) Login(ctx context.Context, s Dto.UserLoginData) DomainLevel.LoginApplicationOutComingData {
	usersData := sa.ReaderDatabase.LoginData(ctx, s.Email)
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

	var BytesID [4]byte
	binary.BigEndian.PutUint32(BytesID[:], uint32(usersData.Id))
	refreshTokenMaker, err := sa.GenereteTokens1.Make()
	if err != nil {
		return DomainLevel.LoginApplicationOutComingData{
			Err: err,
		}
	}
	refreshToken, err := refreshTokenMaker.GetAuthToken(BytesID[:])
	if err != nil {
		return DomainLevel.LoginApplicationOutComingData{
			Jwt: "",
			Rft: "",
			Err: err,
		}
	}

	jwtTokenMaker, err := sa.GenereteToken2.Make()
	if err != nil {
		return DomainLevel.LoginApplicationOutComingData{
			Jwt: "",
			Rft: "",
			Err: err,
		}
	}

	jwtToken, err := jwtTokenMaker.GetAuthToken(BytesID[:])
	if err != nil {
		return DomainLevel.LoginApplicationOutComingData{
			Jwt: "",
			Rft: "",
			Err: err,
		}
	}
	return DomainLevel.LoginApplicationOutComingData{
		Jwt: string(jwtToken),
		Rft: string(refreshToken),
	}
}
