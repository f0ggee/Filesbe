package DomainLevel

import (
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

const (
	FileUrlName = "name"
	TypeFile    = "bool"
)

type RegisterAnswer struct {
	StatusOfOperation string `json:"status_of_operation"`
	UrlToRedirect     string `json:"url_to_redirect"`
	Error             string `json:"error"`
}
type UserCheckAnswer struct {
	UrlToRedirect string `json:"url_to_redirect"`
	Error         string `json:"error"`
}
type AnswerLogin struct {
	StatusOfOperation string `json:"status_of_operation"`
	UrlToRedirect     string `json:"url_to_redirect"`
	ErrorMessage      string `json:"error_message"`
}
type AnswerUrlBuilder struct {
	StatusOperation string `json:"StatusOperation"`
	Url             string `json:"Url"`
	ErrorMessage    string `json:"ErrorMessage"`
}
type AnswerUploaderFileNoEncrypt struct {
	StatusOperation string `json:"StatusOperation"`
	UrlToRedirect   string `json:"UrlRedict"`
	Error           string `json:"Error"`
}

type AuthCheckIncomingData struct {
	Jwt string
	Rft string
}

type OutComingAuthData struct {
	NewJwt          string
	IsNewJwtCreated bool
	Err             error
}
type Auth interface {
	CheckAuthTokens(AuthCheckIncomingData) OutComingAuthData
}
