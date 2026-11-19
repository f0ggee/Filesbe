//This file contains functions that can work with users' sessions (save or modify it)

package httpController

import (
	"Kaban/internal/Drivers/Tokens"
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
	key := os.Getenv("session_key")
	keyIntoBytes, err := hex.DecodeString(key)
	if err != nil {
		panic(err)
	}
	store = sessions.NewCookieStore(keyIntoBytes)
	store.Options = &sessions.Options{
		Path:     "",
		MaxAge:   86400 * 14,
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

type IncomingSessionData struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Jwt     string
	Rt      string
}

func (s SessionConnect) GetSessionData(data IncomingSessionData) ReturnedSessionKey {
	connect, err := store.Get(data.Request, TokenName)
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
	connect, err := store.Get(data.Request, TokenName)
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

type testSession struct {
}

func newTestSession() *testSession {
	return &testSession{}
}

func (t testSession) GetSessionData(data IncomingSessionData) ReturnedSessionKey {

	rf, _ := Tokens.NewRfToken().Make()
	jwt, _ := Tokens.NewJwtToken().Make()
	jwtToken, _ := jwt.GetAuthToken(nil)
	rfToken, _ := rf.GetAuthToken(nil)

	return ReturnedSessionKey{
		Rft:   string(rfToken),
		Jwt:   string(jwtToken),
		Error: nil,
	}

}

func (t testSession) SetNewSession(data IncomingSessionData) ReturnedSessionKey {
	return ReturnedSessionKey{}
}

type testBadSession struct {
}

func newTestBadSession() *testBadSession {
	return &testBadSession{}
}

func (t testBadSession) GetSessionData(data IncomingSessionData) ReturnedSessionKey {
	rf, _ := Tokens.NewTestTokensBadTime().Make()
	jwt, _ := Tokens.NewTestTokensBadTime().Make()

	rfToken, _ := rf.GetAuthToken(nil)
	jwtToken, _ := jwt.GetAuthToken(nil)
	return ReturnedSessionKey{
		Rft:   string(rfToken),
		Jwt:   string(jwtToken),
		Error: nil,
	}
}

func (t testBadSession) SetNewSession(data IncomingSessionData) ReturnedSessionKey {
	return ReturnedSessionKey{}
}
