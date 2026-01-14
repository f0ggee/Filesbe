package Controller

import (
	"Kaban/iternal/Dto"
	Handlers2 "Kaban/iternal/Service/Handlers"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
)

func Auth(Rt string, JwtToken string, session *sessions.Session) (int, string, error, bool) {

	ok := CheckDenyList(Rt)
	if !ok {
		return 0, "", errors.New("token in deny list"), false
	}
	AllowLists := AllowList
	DenyLists := DenyList

	token, err := Handlers2.CheckLifeJWT(JwtToken)
	if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
		err := Handlers2.CheckRefreshTokenLifeTime(Rt)
		if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
			return 0, "", errors.New("refresh token expired"), false
		}
		if err != nil {
			delete(AllowLists, Rt)
			DenyLists[Rt] = time.Now()
			slog.Error("Refresh have been change", err)
			return 0, "", err, false
		}

		NewJwt, NewRet, Id, scrypt, err := Handlers2.GenerateNewTokens(Rt)
		if err != nil {
			slog.Error("error generate new tokens", err)
			return 0, "", err, false
		}
		delete(AllowLists, Rt)
		DenyLists[Rt] = time.Now()
		AllowLists[NewRet] = time.Now()
		session.Values["JT"] = NewJwt
		session.Values["RT"] = NewRet
		return Id, scrypt, nil, false
	}
	if err != nil {
		slog.Error("jwt have baen to change or something more", err)
		return 0, "", err, false
	}

	sa := token.Claims.(*Dto.MyCustomCookie)

	return sa.UserID, sa.Sc, nil, true

}
