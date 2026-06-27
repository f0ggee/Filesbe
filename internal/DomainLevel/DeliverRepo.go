package DomainLevel

import (
	"net/http"
	"time"
)

var StatusResponse struct {
}

const JwtCookieName = "JWTCookie"
const ContentType = "Content-Type"
const RTCookieName = "RTCookie"
const NotStart = "NotStart"
const CookieTimeLive = 1000 * time.Hour
const Break = "BREAK"
const Json = "application/json"
const Success = "Success"
const DomainName = "https://filesbes.com/"
const LocalHostName = "http://localhost:8080/"
const Bots = "Bot"
const RequestId = "RequestId"
const TokenName = "token6"
const InfoPageUrl = "/informationPage"

type ReturnedSessionKey struct {
	Rft   string
	Jwt   string
	Error error
}
type IncomingSessionData struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Jwt     string
	Rt      string
}

type Session interface {
	SetNewSession(data IncomingSessionData) *ReturnedSessionKey
	GetSessionData(data IncomingSessionData) *ReturnedSessionKey
}
type Parses interface {
	JsonParsers(any, *http.Request) error
}
type Register interface {
	ValidateData() error
}
