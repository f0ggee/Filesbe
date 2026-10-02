package Tokens

import (
	"Kaban/internal/DomainLevel"
	"crypto/rand"
	"log/slog"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

type RfToken struct {
	sig jose.Signer
}

func GetNewRfToken() *RfToken {
	return &RfToken{}
}

func (j *RfToken) GetAuthToken(UsefulData []byte) ([]byte, error) {
	cl := jwt.Claims{
		Issuer:   "Kaban",
		Subject:  "u",
		Audience: jwt.Audience{string(UsefulData)},
		Expiry:   jwt.NewNumericDate(time.Now().Add(32 * time.Hour)),
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

func (j *RfToken) IsTokenCorrect(Token []byte) error {

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

func (j *RfToken) SetAdditionalData(bytes []byte) DomainLevel.AuthMaker {
	return j
}
func (j *RfToken) Make() (DomainLevel.Auth, error) {
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
