package ValidatingTokens

import (
	"Kaban/internal/Dto"
	Generator "Kaban/internal/InfrastructureLayer/AuthTokensManage/Creating"
	"crypto/rand"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestChecking_CheckingDenyList(t *testing.T) {
	type args struct {
		s string
	}

	NewRtToken, err := G()
	if err != nil {
		panic(err)
	}
	Dto.DenyList[NewRtToken] = time.Now()
	c := Checking{}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test_1",
			args: args{
				func() string {
					newRtToken, _ := G()
					return newRtToken
				}(),
			},
			want: false,
		},
		{
			name: "Test_2",

			args: args{
				NewRtToken,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.CheckingDenyList(tt.args.s); got != tt.want {
				t.Errorf("CheckingDenyList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func G() (string, error) {
	sax := Generator.CreatingTokens{}

	NewRtToken, err := sax.GenerateRT(Dto.JwtCustomStruct{
		UserID: 12,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "",
			Subject:   "",
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			ID:        rand.Text(),
		},
	})
	return NewRtToken, err
}
