package Tokens

import (
	"Kaban/internal/DomainLevel"
	"crypto/rand"
	"log/slog"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

type JwtToken struct {
	sig jose.Signer
}

func GetNewJwtToken() JwtToken {
	return JwtToken{}
}

func (j *JwtToken) GetAuthToken(UsefulData []byte) ([]byte, error) {
	cl := jwt.Claims{
		Issuer:   "Kaban",
		Subject:  "u",
		Audience: jwt.Audience{string(UsefulData)},
		Expiry:   jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ID:       rand.Text(),
	}
	raw, err := jwt.Signed(j.sig).Claims(cl).Serialize()
	if err != nil {
		slog.Error("GetAuthToken jwt: error to serialize", "ERROR", err)
		return nil, ErrorMakeToken
	}
	return []byte(raw), nil
}

func (j *JwtToken) IsTokenCorrect(Token []byte) error {

	parsedToken, err := jwt.ParseSigned(string(Token), []jose.SignatureAlgorithm{jose.HS512})
	if err != nil {
		slog.Error("JwtToken IsTokenCorrect: an error happened during parsing", "ERROR", err)
		return ErrorParseToken
	}

	claims := jwt.Claims{}

	if err := parsedToken.Claims(Key, &claims); err != nil {
		slog.Error("JwtToken IsTokenCorrect: an error happened")
		return ErrorTokenAuthFail
	}

	if err := claims.Validate(jwt.Expected{
		Time: time.Now(),
	}); err != nil {
		return ErrorTokenFail
	}
	return nil
}

func (j *JwtToken) SetAdditionalData(bytes []byte) DomainLevel.AuthMaker {
	return j
}
func (j *JwtToken) Make() (DomainLevel.Auth, error) {
	d, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.HS512,
		Key:       Key,
	}, nil)

	if err != nil {
		slog.Error("Make JwtKey: error to produce a new sginer", "ERROR", err)
		return nil, ErrorTokenKey
	}

	j.sig = d
	return j, nil
}
