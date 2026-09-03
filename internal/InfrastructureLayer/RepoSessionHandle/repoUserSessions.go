package RepoSessionHandle

import (
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

type NewSessionConnect struct {
	Activity *sessions.CookieStore
}

func GetNewSessionConnect(activity *sessions.CookieStore) *NewSessionConnect {
	return &NewSessionConnect{Activity: activity}
}

type IncomingSessionData struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Jwt     string
	Rt      string
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

func (s *NewSessionConnect) getUserConnect(r *http.Request) (*sessions.Session, error) {
	return s.Activity.Get(r, TokenName)
}

func GetCookieStore() *sessions.CookieStore {
	var store1z, err = hex.DecodeString(os.Getenv("KEY1"))
	if err != nil {
		slog.Error("SessionControl: error to create a cookie", "Err", err)
		panic(err)
	}
	Store := sessions.NewCookieStore(store1z)
	return Store
}
func (s NewSessionConnect) GetSessionData(data IncomingSessionData) ReturnedSessionKey {
	connect, err := s.getUserConnect(data.Request)
	if err != nil {
		slog.Error("GetSessionData: error to get an active connect", "error", err)
		return ReturnedSessionKey{
			Error: errors.New(ErrorGetCookie),
		}
	}
	if connect.Options.MaxAge == 0 {
		return ReturnedSessionKey{Error: errors.New(ErrorAuthExpired)}
	}
	rtToken, isKeyExist := connect.Values[RTCookieName].(string)
	if !isKeyExist {
		return ReturnedSessionKey{Error: errors.New(ErrorAuthExpired)}
	}
	jwts, isKeyExist := connect.Values[JwtCookieName].(string)
	if !isKeyExist {
		return ReturnedSessionKey{Error: errors.New(ErrorAuthExpired)}
	}
	return ReturnedSessionKey{
		Rft: rtToken,
		Jwt: jwts,
	}
}
func (s *NewSessionConnect) SetNewSession(data IncomingSessionData) ReturnedSessionKey {
	connect, err := s.getUserConnect(data.Request)
	if err != nil {
		return ReturnedSessionKey{Error: errors.New(ErrorGetCookie)}
	}
	if connect.Options.MaxAge == 0 {
		return ReturnedSessionKey{Error: errors.New(ErrorAuthExpired)}
	}

	if data.Jwt != "" {
		connect.Values[JwtCookieName] = data.Jwt
	}

	if data.Rt != "" {
		connect.Values[RTCookieName] = data.Rt
	}
	connect.Options = &sessions.Options{
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(CookieTimeLive),
	}

	if err := connect.Save(data.Request, data.Writer); err != nil {
		slog.Error("Error in save cookie", "Err", err)
		return ReturnedSessionKey{Error: errors.New(ErrorSaveCookie)}
	}
	return ReturnedSessionKey{
		Rft:   "",
		Jwt:   "",
		Error: nil,
	}
}

type Mocks struct {
}

func (m Mocks) GetSessionData(data IncomingSessionData) ReturnedSessionKey {

	return ReturnedSessionKey{}
}

func (m Mocks) SetNewSession(data IncomingSessionData) ReturnedSessionKey {
	return ReturnedSessionKey{}
}
