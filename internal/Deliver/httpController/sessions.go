//This file contains functions that can work with users' sessions(save it or modify it)

package httpController

import (
	"encoding/hex"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

const RTCookieName = "RTCookie"
const JwtCookieName = "JWTCookie"
const TokenName = "token6"

var store *sessions.CookieStore

func init() {
	key := os.Getenv("KEY_FOR_JWT")
	keyIntoBytes, err := hex.DecodeString(key)
	if err != nil {
		panic(err)
	}
	store = sessions.NewCookieStore(keyIntoBytes)
	store.Options = &sessions.Options{
		Path:     "",
		MaxAge:   1000 * 3600,
		Secure:   os.Getenv("app_product") == "product",
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}
}

type Session interface {
	GetSessionData(data IncomingSessionData) ReturnedSessionKey
	SetNewSession(data IncomingSessionData) ReturnedSessionKey
}
type ReturnedSessionKey struct {
	Rft   string
	Jwt   string
	Error error
}

type SessionConnect struct{}

func GetNewSessionConnect() SessionConnect {
	return SessionConnect{}
}

func (s SessionConnect) getUserConnect(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, TokenName)
}

type IncomingSessionData struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Jwt     string
	Rt      string
}

func (s SessionConnect) GetSessionData(data IncomingSessionData) ReturnedSessionKey {
	connect, err := s.getUserConnect(data.Request)
	if err != nil {
		slog.Error("GetSessionData: error to get an active connect", "error", err)
		return ReturnedSessionKey{
			Error: ErrorGetCookie,
		}
	}
	if connect.Options.MaxAge == 0 {
		return ReturnedSessionKey{Error: ErrorAuthExpired}
	}
	rtToken, isKeyExist := connect.Values[RTCookieName].(string)
	if !isKeyExist {
		return ReturnedSessionKey{Error: ErrorAuthExpired}
	}
	jwts, isKeyExist := connect.Values[JwtCookieName].(string)
	if !isKeyExist {
		return ReturnedSessionKey{Error: ErrorAuthExpired}
	}
	return ReturnedSessionKey{
		Rft: rtToken,
		Jwt: jwts,
	}
}
func (s SessionConnect) SetNewSession(data IncomingSessionData) ReturnedSessionKey {
	connect, err := s.getUserConnect(data.Request)
	if err != nil {
		return ReturnedSessionKey{Error: ErrorGetCookie}
	}
	if connect.Options.MaxAge == 0 {
		return ReturnedSessionKey{Error: ErrorAuthExpired}
	}

	if data.Jwt != "" {
		connect.Values[JwtCookieName] = data.Jwt
	}
	if data.Rt != "" {
		connect.Values[RTCookieName] = data.Rt
	}
	if err = store.Save(data.Request, data.Writer, connect); err != nil {
		slog.Error("SetNewSession: error to save data", "ERROR", err)
		return ReturnedSessionKey{Error: ErrorSaveCookie}
	}
	return ReturnedSessionKey{
		Rft:   "",
		Jwt:   "",
		Error: nil,
	}
}

type Mocks struct {
}

func GetNewMocks() *Mocks {
	return &Mocks{}
}

func (m Mocks) GetSessionData(data IncomingSessionData) ReturnedSessionKey {

	return ReturnedSessionKey{}
}

func (m Mocks) SetNewSession(data IncomingSessionData) ReturnedSessionKey {
	return ReturnedSessionKey{}
}
