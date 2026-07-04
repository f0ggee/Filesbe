package RepoSessionHandle

import (
	"Kaban/internal/DomainLevel"
	"errors"
	"log/slog"
)

func (s SessionConnect) GetSessionData(data IncomingSessionData) ReturnedSessionKey {
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
