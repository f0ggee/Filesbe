package Handlers

import (
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage/Creating"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage/ValidatingTokens"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Test_funcName(t *testing.T) {
	type args struct {
		Rt       string
		JwtToken string
		sa       *HandlerPackAuthTokens
	}

	sa := &HandlerPackAuthTokens{
		Checking:        ValidatingTokens.Checking{},
		GeneratingToken: Creating.CreatingTokens{},
	}
	tests := []struct {
		name    string
		Descr   string
		args    args
		wantErr bool
	}{
		{
			name: "Test_1",
			args: args{

				Rt:       R(time.Now().Add(-10000*time.Hour), &Creating.CreatingTokens{}),
				JwtToken: J(time.Now().Add(-10000*time.Hour), &Creating.CreatingTokens{}),
				sa:       sa,
			},
			wantErr: true,
		},
		{
			name: "Test_2",

			args: args{
				JwtToken: "",
				Rt:       R(time.Now().Add(10000*time.Hour), &Creating.CreatingTokens{}),
				sa:       sa,
			},
			wantErr: false,
		},
		{
			name: "Test_3",
			args: args{
				Rt:       "",
				JwtToken: "",
				sa:       sa,
			},
			wantErr: true,
		},
		{
			name: "Test_4",
			args: args{
				JwtToken: J(time.Now().Add(10000*time.Hour), &Creating.CreatingTokens{}),
				Rt:       "",
				sa:       sa,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := funcName(tt.args.Rt, tt.args.JwtToken, tt.args.sa); (err != nil) != tt.wantErr {
				t.Errorf("funcName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func R(time2 time.Time, sa *Creating.CreatingTokens) string {

	e, err := sa.GenerateRT(Dto.JwtCustomStruct{
		UserID: 0,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time2),
		},
	})
	if err != nil {
		panic(err)
	}
	return e

}

func J(time2 time.Time, sa *Creating.CreatingTokens) string {
	e1, err := sa.GenerateJWT(Dto.JwtCustomStruct{
		UserID: 0,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time2),
		},
	})
	if err != nil {
		panic(err)
	}
	return e1
}

func TestHandlerPackAuthTokens_AuthTest(t *testing.T) {
	type fields struct {
		Sa *HandlerPackAuthTokens
	}
	sa := &HandlerPackAuthTokens{
		Checking:        ValidatingTokens.Checking{},
		GeneratingToken: Creating.CreatingTokens{},
	}
	type args struct {
		Rt       string
		JwtToken string
	}

	durationPosi := 10000 * time.Hour
	durationNegative := -10000 * time.Hour
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test_1",
			fields: fields{
				Sa: sa,
			},
			args: args{
				Rt:       R(time.Now().Add(durationPosi), &Creating.CreatingTokens{}),
				JwtToken: J(time.Now().Add(durationPosi), &Creating.CreatingTokens{}),
			},

			wantErr: false,
		},
		{
			name: "Test_2",
			fields: fields{
				Sa: sa,
			},
			args: args{
				JwtToken: "",
				Rt:       R(time.Now().Add(durationPosi), &Creating.CreatingTokens{}),
			},
			wantErr: false,
		},
		{
			name: "Test_3",
			fields: fields{
				Sa: sa,
			},
			args: args{
				JwtToken: "",
				Rt:       "",
			},
			wantErr: true,
		},
		{
			name: "Test_4",
			fields: fields{
				Sa: sa,
			},
			args: args{
				JwtToken: J(time.Now().Add(durationNegative), &Creating.CreatingTokens{}),
				Rt:       R(time.Now().Add(durationNegative), &Creating.CreatingTokens{}),
			},
			wantErr: true,
		},
		{
			name: "Test_5",
			fields: fields{
				Sa: sa,
			},
			args: args{
				JwtToken: J(time.Now().Add(durationNegative), &Creating.CreatingTokens{}),
				Rt:       R(time.Now().Add(durationPosi), &Creating.CreatingTokens{}),
			},

			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sa := &HandlerPackAuthTokens{
				Manage:          tt.fields.Sa.Manage,
				GeneratingToken: tt.fields.Sa.GeneratingToken,
				Checking:        tt.fields.Sa.Checking,
			}
			_, err := sa.AuthTest(tt.args.Rt, tt.args.JwtToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthTest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
