package Tokens

import (
	"Kaban/internal/DomainLevel"
	"crypto/rand"
	"io"
	"log/slog"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

type TestTokensBadTime struct {
	sig jose.Signer
}

func NewTestTokensBadTime() *TestTokensBadTime {
	return &TestTokensBadTime{}
}

func (t TestTokensBadTime) GetAuthToken(UsefulData []byte) ([]byte, error) {
	cl := jwt.Claims{
		Issuer:   "Kaban",
		Subject:  "u",
		Audience: jwt.Audience{string(UsefulData)},
		Expiry:   jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ID:       rand.Text(),
	}
	raw, err := jwt.Signed(t.sig).Claims(cl).Serialize()
	if err != nil {
		slog.Error("GetAuthToken jwt: error to serialize", "ERROR", err)
		return nil, ErrorMakeToken
	}
	return []byte(raw), nil
}

func (t TestTokensBadTime) IsTokenCorrect(Token []byte) error {
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
	err = claims.Validate(jwt.Expected{
		Time: time.Now(),
	})

	if claims.Expiry != nil {
		return ErrorTokenExpired
	}
	if err != nil {
		return ErrorTokenFail
	}
	return nil
}

func (t TestTokensBadTime) SetAdditionalData(bytes []byte) DomainLevel.AuthMaker {
	return nil
}

func (t TestTokensBadTime) Make() (DomainLevel.Auth, error) {
	d, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.HS512,
		Key:       Key,
	}, nil)
	if err != nil {
		slog.Error("Make JwtKey: error to produce a new sginer", "ERROR", err)
		return nil, ErrorTokenKey
	}
	t.sig = d
	return t, nil
}

type testTokensBadKey struct {
	sig jose.Signer
}

func newTestTokensBadKey() *testTokensBadKey {
	return &testTokensBadKey{}
}

func (t *testTokensBadKey) GetAuthToken(UsefulData []byte) ([]byte, error) {
	cl := jwt.Claims{
		Issuer:   "Kaban",
		Subject:  "u",
		Audience: jwt.Audience{string(UsefulData)},
		Expiry:   jwt.NewNumericDate(time.Now().Add(32 * time.Hour)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ID:       rand.Text(),
	}
	raw, err := jwt.Signed(t.sig).Claims(cl).Serialize()
	if err != nil {
		slog.Error("GetAuthToken jwt: error to serialize", "ERROR", err)
		return nil, ErrorMakeToken
	}
	return []byte(raw), nil
}

func (t *testTokensBadKey) IsTokenCorrect(Token []byte) error {

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

	if claims.Expiry != nil {
		return ErrorTokenExpired
	}
	if err := claims.Validate(jwt.Expected{
		Time: time.Now(),
	}); err != nil {
		return ErrorTokenFail
	}
	return nil

}

func (t *testTokensBadKey) SetAdditionalData(bytes []byte) DomainLevel.AuthMaker {
	return t
}

func (t *testTokensBadKey) Make() (DomainLevel.Auth, error) {
	newKey := make([]byte, 64)
	io.ReadFull(rand.Reader, newKey)
	d, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.HS512,
		Key:       newKey,
	}, nil)
	if err != nil {
		slog.Error("Make JwtKey: error to produce a new signer", "ERROR", err)
		return nil, ErrorTokenKey
	}
	t.sig = d
	return t, nil
}

type testTokensBadAlg struct {
	sig jose.Signer
}

func newTestTokensBadAlg() *testTokensBadAlg {
	return &testTokensBadAlg{}
}

func (t *testTokensBadAlg) GetAuthToken(UsefulData []byte) ([]byte, error) {
	cl := jwt.Claims{
		Issuer:   "Kaban",
		Subject:  "u",
		Audience: jwt.Audience{string(UsefulData)},
		Expiry:   jwt.NewNumericDate(time.Now().Add(32 * time.Hour)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ID:       rand.Text(),
	}
	raw, err := jwt.Signed(t.sig).Claims(cl).Serialize()
	if err != nil {
		slog.Error("GetAuthToken jwt: error to serialize", "ERROR", err)
		return nil, ErrorMakeToken
	}
	return []byte(raw), nil

}

func (t *testTokensBadAlg) IsTokenCorrect(Token []byte) error {
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

	if claims.Expiry != nil {
		return ErrorTokenExpired
	}
	if err := claims.Validate(jwt.Expected{
		Time: time.Now(),
	}); err != nil {
		return ErrorTokenFail
	}
	return nil
}

func (t *testTokensBadAlg) SetAdditionalData(bytes []byte) DomainLevel.AuthMaker {
	return nil
}

func (t *testTokensBadAlg) Make() (DomainLevel.Auth, error) {

	d, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.HS256,
		Key:       Key,
	}, nil)
	if err != nil {
		slog.Error("Make JwtKey: error to produce a new sginer", "ERROR", err)
		return nil, ErrorTokenKey
	}
	t.sig = d
	return t, nil

}
