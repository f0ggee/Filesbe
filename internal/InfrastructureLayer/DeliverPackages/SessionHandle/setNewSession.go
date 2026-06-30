package SessionHandle

import (
	"Kaban/internal/DomainLevel"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
)

func (s *SessionConnect) SetNewSession(data DomainLevel.IncomingSessionData) DomainLevel.ReturnedSessionKey {
	connect, err := s.getUserConnect(data.Request)
	if err != nil {
		return DomainLevel.ReturnedSessionKey{Error: errors.New(DomainLevel.ErrorGetCookie)}
	}
	if connect.Options.MaxAge == 0 {
		return DomainLevel.ReturnedSessionKey{Error: errors.New(DomainLevel.ErrorAuthExpired)}
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

	s.mut.Lock()
	if err := connect.Save(data.Request, data.Writer); err != nil {
		slog.Error("Error in save cookie", "Err", err)
		return DomainLevel.ReturnedSessionKey{Error: errors.New(DomainLevel.ErrorSaveCookie)}

	}
	s.mut.Unlock()
	return DomainLevel.ReturnedSessionKey{
		Rft:   "",
		Jwt:   "",
		Error: nil,
	}
}
