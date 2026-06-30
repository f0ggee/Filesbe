package DomainLevel

import (
	"io"
	"net/http"
	"time"
)

const JwtCookieName = "JWTCookie"
const ContentType = "Content-Type"
const RTCookieName = "RTCookie"

const CookieTimeLive = 1000 * time.Hour

const Json = "application/json"
const Success = "Success"
const DomainName = "https://filesbes.com/"
const LocalHostName = "http://localhost:8080/"
const Bots = "Bot"
const RequestId = "RequestId"
const TokenName = "token6"
const InfoPageUrl = "/informationPage"

type RegisterAnswer struct {
	StatusOfOperation string `json:"status_of_operation"`
	UrlToRedirect     string `json:"url_to_redirect"`
	Error             string `json:"error"`
}
type UserCheckAnswer struct {
	UrlToRedirect string `json:"url_to_redirect"`
	Error         string `json:"error"`
}

type RegisterErrorIncomingData struct {
	W         http.ResponseWriter
	Error     error
	Operation error
}
type RegisterGoodIncomingData struct {
	W         http.ResponseWriter
	Operation string
	Redirect  string
}

type UserCheckBadIncomingData struct {
	W        http.ResponseWriter
	Redirect string
	Err      error
}
type Session interface {
	SetNewSession(data IncomingSessionData) ReturnedSessionKey
	GetSessionData(data IncomingSessionData) *ReturnedSessionKey
}
type Parses interface {
	JsonParsers(any, io.ReadCloser) error
}
type Register interface {
	ErrorAnswer(RegisterErrorIncomingData)
	GoodAnswer(RegisterGoodIncomingData)
}

type UsersCheck interface {
	BadAnswer(UserCheckBadIncomingData)
}
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
