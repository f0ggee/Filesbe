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

func TestNewRegisterApplication_RegisterApplication(t *testing.T) {
	type fields struct {
		CheckingDb DomainLevel.CheckingDb
		WriterDb   DomainLevel.WriteDb
		Generator  DomainLevel.CryptoGenerating
		Rft        DomainLevel.AuthMaker
		Jwt        DomainLevel.AuthMaker
	}
	type args struct {
		ctx context.Context
		de  *Dto.UserDataRegister
	}

	defaultFields := fields{
		CheckingDb: DatabaseControl.NewCheckerTest(),
		WriterDb:   DatabaseControl.NewWriterTest(),
		Generator:  Crypto.GetNewGenerating(),
		Rft:        Tokens.NewJwtToken(),
		Jwt:        Tokens.NewRfToken(),
	}

	defaultParams := args{
		ctx: context.Background(),
		de: &Dto.UserDataRegister{
			Name:     "1234",
			Email:    "33",
			Password: []byte("Hello!"),
		},
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
				CheckingDb: defaultFields.CheckingDb,
				WriterDb:   DatabaseControl.NewWriterError(),
				Generator:  defaultFields.Generator,
				Rft:        defaultFields.Jwt,
				Jwt:        defaultFields.Rft,
			},
			args:      defaultParams,
			wantError: true,
		},
		{
			name: "test3",
			fields: fields{
				CheckingDb: defaultFields.CheckingDb,
				WriterDb:   DatabaseControl.NewWriterError(),
				Generator:  defaultFields.Generator,
				Rft:        defaultFields.Rft,
				Jwt:        defaultFields.Jwt,
			},
			args:      defaultParams,
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sa := &NewRegisterApplication{
				CheckingDb: tt.fields.CheckingDb,
				WriterDb:   tt.fields.WriterDb,
				Generator:  tt.fields.Generator,
				Rft:        tt.fields.Rft,
				Jwt:        tt.fields.Jwt,
			}
			if got := sa.RegisterApplication(tt.args.ctx, tt.args.de); (got.Err != nil) != tt.wantError {
				t.Errorf("RegisterApplication() = %v, want %v", got, tt.wantError)
			}
			for i := range tt.args.de.Password {
				if tt.args.de.Password[i] != 0 {
					t.Errorf("data isn't deleted ")
				}
			}

		})
	}
}
