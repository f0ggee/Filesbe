package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/Crypto"
	"Kaban/internal/InfrastructureLayer/DatabaseControl"
	"Kaban/internal/InfrastructureLayer/Tokens"
	"context"
	"testing"
)

func TestNewLogin_LoginApplication(t *testing.T) {
	type fields struct {
		ReaderDatabase DomainLevel.ReadDb
		Validate       DomainLevel.CryptoValidating
		Rf             DomainLevel.AuthMaker
		Jwt            DomainLevel.AuthMaker
	}
	type args struct {
		ctx context.Context
		s   Dto.UserLoginData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test1",
			fields: fields{
				ReaderDatabase: DatabaseControl.NewReaderTest(),
				Validate:       Crypto.NewCheckingTest(),
				Rf:             Tokens.NewJwtToken(),
				Jwt:            Tokens.NewRfToken(),
			},
			args: args{
				ctx: context.WithValue(context.Background(), "mode", 1),
				s: Dto.UserLoginData{
					Email:    "testmail@gmail.com",
					Password: "12345678",
				},
			},
			want: false,
		},
		{
			name: "test2",
			fields: fields{
				ReaderDatabase: DatabaseControl.NewReaderTest(),
				Validate:       Crypto.NewCheckingTestError(),
				Rf:             Tokens.NewRfToken(),
				Jwt:            Tokens.NewJwtToken(),
			},
			args: args{
				ctx: context.WithValue(context.Background(), "mode", 1),
				s: Dto.UserLoginData{
					Email:    "passw@gmail.com",
					Password: "12345678",
				},
			},
			want: true,
		},
		{
			name: "test3",
			fields: fields{
				ReaderDatabase: DatabaseControl.NewReaderTest(),
				Validate:       Crypto.GetNewChecking(),
				Rf:             Tokens.NewJwtToken(),
				Jwt:            Tokens.NewRfToken(),
			},
			args: args{
				ctx: context.WithValue(context.Background(), "mode", 2),
				s: Dto.UserLoginData{
					Email:    "psjajdh@gmail.com ",
					Password: "123467",
				},
			},
			want: true,
		},
		{
			name: "test4",
			fields: fields{
				ReaderDatabase: DatabaseControl.NewReaderTest(),
				Validate:       Crypto.GetNewChecking(),
				Rf:             Tokens.NewJwtToken(),
				Jwt:            Tokens.NewRfToken(),
			},
			args: args{
				ctx: context.WithValue(context.Background(), "mode", 3),
				s: Dto.UserLoginData{
					Email:    "psasd@gmail.com",
					Password: "1234",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sa := &NewLogin{
				ReaderDatabase: tt.fields.ReaderDatabase,
				Validate:       tt.fields.Validate,
				Rf:             tt.fields.Rf,
				Jwt:            tt.fields.Jwt,
			}
			if got := sa.LoginApplication(tt.args.ctx, tt.args.s); !tt.want && got.Err != nil {
				t.Errorf("LoginApplication() = %v, want %v", got, tt.want)
			}
		})
	}
}
