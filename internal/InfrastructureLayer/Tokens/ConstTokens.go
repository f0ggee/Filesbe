package Tokens

import "os"

var (
	Key []byte
)

const ErrorTokenKey = "the token key wasn't found"
const ErrorMakeToken = "the token can't be created"
const ErrorParseToken = "can't parse a token"
const ErrorTokenAuthFail = "token's auth is failed"
const ErrorTokenFail = "a token isn't correct"

func init() {
	k := os.Getenv("KEY_FOR_JWT")
	if k == "" {
		panic(ErrorTokenKey)
	}
	Key = []byte(k)
}
