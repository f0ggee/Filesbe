package Tokens

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var (
	Key []byte
)

var ErrorTokenKey = errors.New("the token key wasn't found")
var ErrorMakeToken = errors.New("the token can't be created")
var ErrorParseToken = errors.New("can't parse a token")
var ErrorTokenAuthFail = errors.New("token's auth is failed")
var ErrorTokenFail = errors.New("a token isn't correct")

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	k := os.Getenv("KEY_FOR_JWT")
	if k == "" {
		panic(ErrorTokenKey)
	}
	Key = []byte(k)
}
