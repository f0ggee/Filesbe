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

type NewLogin struct {
	ReaderDatabase DomainLevel.ReadDb
	Validate       DomainLevel.CryptoValidating
	Rf             DomainLevel.AuthMaker
	Jwt            DomainLevel.AuthMaker
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
	refreshTokenMaker, err := sa.Rf.Make()
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

	jwtTokenMaker, err := sa.Jwt.Make()
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
