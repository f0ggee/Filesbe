package ValidatingTokens

import (
	"Kaban/internal/Dto"
	Generator "Kaban/internal/InfrastructureLayer/AuthTokensManage/Creating"
	c "crypto/rand"
	"io"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestChecking_CheckJwt(t *testing.T) {
	type args struct {
		JWT string
	}
	g := Generator.CreatingTokens{}
	tests := []struct {
		name    string
		descr   string
		args    args
		wantErr bool
	}{
		{
			name:  "Test_1",
			descr: "A test with a valid JWT",
			args: args{
				func() string {
					NewFineToken, err := g.GenerateJWT(Dto.JwtCustomStruct{
						UserID: rand.Int(),
						RegisteredClaims: jwt.RegisteredClaims{
							Issuer:    c.Text(),
							Subject:   c.Text(),
							ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
							IssuedAt:  jwt.NewNumericDate(time.Now()),
							ID:        c.Text(),
						},
					})
					if err != nil {
						panic(err)
					}
					return NewFineToken
				}(),
			},
			wantErr: false,
		},
		{
			name:  "Test_2",
			descr: "A test with a  expired JWT token",
			args: args{
				func() string {

					NewFineToken, err := g.GenerateJWT(Dto.JwtCustomStruct{
						UserID: rand.Int(),
						RegisteredClaims: jwt.RegisteredClaims{
							Issuer:    c.Text(),
							Subject:   c.Text(),
							ExpiresAt: jwt.NewNumericDate(time.Now().Add(-21 * time.Hour)),
							IssuedAt:  jwt.NewNumericDate(time.Now().Add(-22 * time.Hour)),
						},
					})
					if err != nil {
						panic(err)
					}
					return NewFineToken
				}(),
			},
			wantErr: true,
		},
		{
			name:  "Test_3",
			descr: "A test with a nil JWT token",
			args: args{
				JWT: "",
			},
			wantErr: true,
		},
		{
			name:  "Test_4",
			descr: "A test with a Jwt token which was signed by another key",
			args: args{

				GenerateWrongKey(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checking := Checking{}
			t.Logf("Description about a test : %s", tt.descr)
			_, err := checking.CheckJwt(tt.args.JWT)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckJwt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func GenerateWrongKey() string {

	NewKey := make([]byte, 32)
	if _, err := io.ReadFull(c.Reader, NewKey); err != nil {
		panic(err)
	}
	JwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Dto.JwtCustomStruct{})

	signedString, err := JwtToken.SignedString(NewKey)
	if err != nil {
		panic(err)
	}
	return signedString
}
