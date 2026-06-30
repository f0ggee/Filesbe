package SessionHandle

import (
	"Kaban/internal/DomainLevel"
	"errors"
	"log/slog"
)

func (s SessionConnect) GetSessionData(data DomainLevel.IncomingSessionData) *DomainLevel.ReturnedSessionKey {
	connect, err := s.getUserConnect(data.Request)
	if err != nil {
		slog.Error("GetSessionData: error to get an active connect", "error", err)
		return &DomainLevel.ReturnedSessionKey{
			Error: errors.New(DomainLevel.ErrorGetCookie),
		}
	}
	if connect.Options.MaxAge == 0 {
		return &DomainLevel.ReturnedSessionKey{Error: errors.New(DomainLevel.ErrorAuthExpired)}
	}

	s.mut.RLock()
	rtToken, isKeyExist := connect.Values[DomainLevel.RTCookieName].(string)
	if !isKeyExist {
		return &DomainLevel.ReturnedSessionKey{Error: errors.New(DomainLevel.ErrorAuthExpired)}
	}

	s.mut.RUnlock()
	s.mut.RLock()
	jwts, isKeyExist := connect.Values[DomainLevel.JwtCookieName].(string)
	if !isKeyExist {
		return &DomainLevel.ReturnedSessionKey{Error: errors.New(DomainLevel.ErrorAuthExpired)}
	}
	s.mut.RUnlock()

	return &DomainLevel.ReturnedSessionKey{
		Rft: rtToken,
		Jwt: jwts,
	}
}
