package RepoSessionHandle

import "time"

const RTCookieName = "RTCookie"
const JwtCookieName = "JWTCookie"
const CookieTimeLive = 1000 * time.Hour
const TokenName = "token6"

const (
	ErrorUserNotAuthed = "user's auth is expired"
	ErrorUserToken     = "user's auth isn't valid"
)
