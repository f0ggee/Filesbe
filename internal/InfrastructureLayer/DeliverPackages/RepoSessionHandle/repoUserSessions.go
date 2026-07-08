package RepoSessionHandle

import (
	"Kaban/internal/DomainLevel"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/sessions"
)

type NewSessionConnect struct {
	Activity *sessions.CookieStore
	mut      *sync.RWMutex
}

func GetNewSessionConnect(activity *sessions.CookieStore, mut *sync.RWMutex) *NewSessionConnect {
	return &NewSessionConnect{Activity: activity, mut: mut}
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
	s.mut.RLock()
	defer s.mut.RUnlock()
	return s.Activity.Get(r, DomainLevel.TokenName)
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
			Error: errors.New(DomainLevel.ErrorGetCookie),
		}
	}
	if connect.Options.MaxAge == 0 {
		return ReturnedSessionKey{Error: errors.New(DomainLevel.ErrorAuthExpired)}
	}
	rtToken, isKeyExist := connect.Values[DomainLevel.RTCookieName].(string)
	if !isKeyExist {
		return ReturnedSessionKey{Error: errors.New(DomainLevel.ErrorAuthExpired)}
	}
	jwts, isKeyExist := connect.Values[DomainLevel.JwtCookieName].(string)
	if !isKeyExist {
		return ReturnedSessionKey{Error: errors.New(DomainLevel.ErrorAuthExpired)}
	}
	return ReturnedSessionKey{
		Rft: rtToken,
		Jwt: jwts,
	}
}
func (s *NewSessionConnect) SetNewSession(data IncomingSessionData) ReturnedSessionKey {
	connect, err := s.getUserConnect(data.Request)
	if err != nil {
		return ReturnedSessionKey{Error: errors.New(DomainLevel.ErrorGetCookie)}
	}
	if connect.Options.MaxAge == 0 {
		return ReturnedSessionKey{Error: errors.New(DomainLevel.ErrorAuthExpired)}
	}

	if data.Jwt != "" {
		connect.Values[DomainLevel.JwtCookieName] = data.Jwt
	}

	if data.Rt != "" {
		connect.Values[DomainLevel.RTCookieName] = data.Rt
	}
	connect.Options = &sessions.Options{
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(1000 * time.Hour),
	}

	if err := connect.Save(data.Request, data.Writer); err != nil {
		slog.Error("Error in save cookie", "Err", err)
		return ReturnedSessionKey{Error: errors.New(DomainLevel.ErrorSaveCookie)}

	}
	return ReturnedSessionKey{
		Rft:   "",
		Jwt:   "",
		Error: nil,
	}
}
