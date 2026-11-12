package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Drivers/Crypto"
	"Kaban/internal/Drivers/DatabaseControl"
	"Kaban/internal/Drivers/Tokens"
	"Kaban/internal/Dto"
	"context"
	"testing"
)

func TestNewLogin_LoginApplication(t *testing.T) {

	password := []byte("Password")

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
	defaultParams := args{
		ctx: context.Background(),
		s: Dto.UserLoginData{
			Email:    "12345@gmail.com",
			Password: password,
		},
	}

	defaultFields := fields{
		ReaderDatabase: DatabaseControl.NewReaderTest(),
		Validate:       Crypto.GetNewChecking(),
		Rf:             Tokens.NewRfToken(),
		Jwt:            Tokens.NewJwtToken(),
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantError bool
	}{
		{
			name:      "test1",
			fields:    defaultFields,
			args:      defaultParams,
			wantError: false,
		},
		{
			name: "test2",
			fields: fields{
				ReaderDatabase: DatabaseControl.NewReaderErrorTest(),
				Validate:       defaultFields.Validate,
				Rf:             defaultFields.Jwt,
				Jwt:            defaultFields.Rf,
			},
			args:      defaultParams,
			wantError: true,
		},
		{
			name:   "test3",
			fields: defaultFields,
			args: args{
				ctx: defaultParams.ctx,
				s: Dto.UserLoginData{
					Email:    "1234@gmail.com",
					Password: []byte("Data"),
				},
			},
			wantError: true,
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
			if got := sa.LoginApplication(tt.args.ctx, tt.args.s); (got.Err != nil) != tt.wantError {
				t.Errorf("LoginApplication() = %v, want %v", got, tt.wantError)
			}
		})
	}
}
