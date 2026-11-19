package Tokens

import (
	"Kaban/internal/DomainLevel"
	"testing"
)

func TestTokens(t *testing.T) {
	tests := []struct {
		name       string
		authSystem DomainLevel.AuthMaker
		token      DomainLevel.Auth // An implementation will generate a token for checking by authSystem
		wantError  bool
		errorPart  int
	}{

		{
			name:       "test1",
			authSystem: NewRfToken(),
			token: func() DomainLevel.Auth {
				d, _ := NewRfToken().Make()
				return d
			}(),
			wantError: false,
			errorPart: 0,
		},
		{
			name:       "test2",
			authSystem: NewJwtToken(),
			token: func() DomainLevel.Auth {
				d, _ := NewJwtToken().Make()
				return d
			}(),
			wantError: false,
			errorPart: 0,
		},
		{
			name:       "test3",
			authSystem: NewJwtToken(),
			token: func() DomainLevel.Auth {
				d, _ := NewTestTokensBadTime().Make()
				return d
			}(),
			wantError: true,
			errorPart: 3,
		},
		{
			name:       "test4",
			authSystem: NewRfToken(),
			token: func() DomainLevel.Auth {
				d, _ := NewTestTokensBadTime().Make()
				return d
			}(),
			wantError: true,
			errorPart: 3,
		},
		{
			name:       "test5",
			authSystem: NewRfToken(),
			token: func() DomainLevel.Auth {
				d, _ := newTestTokensBadKey().Make()
				return d
			}(),
			wantError: true,
			errorPart: 3,
		},
		{
			name:       "test6",
			authSystem: NewJwtToken(),
			token: func() DomainLevel.Auth {
				d, _ := newTestTokensBadKey().Make()
				return d
			}(),
			wantError: true,
			errorPart: 3,
		},
		{
			name:       "test7",
			authSystem: NewRfToken(),
			token: func() DomainLevel.Auth {
				d, _ := newTestTokensBadAlg().Make()
				return d
			}(),
			wantError: true,
			errorPart: 3,
		},
		{
			name:       "test8",
			authSystem: NewJwtToken(),
			token: func() DomainLevel.Auth {
				d, _ := newTestTokensBadAlg().Make()
				return d
			}(),
			wantError: true,
			errorPart: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkedSystem, err := tt.authSystem.Make()
			if (err != nil) != tt.wantError && tt.errorPart == 1 {
				t.Errorf("TestTokens() Make wanted = %v  but got = %v ", tt.wantError, err)
				return
			}
			token, err := tt.token.GetAuthToken(nil)
			if (err != nil) != tt.wantError && tt.errorPart == 2 {
				t.Errorf("TestTokens() GetAuthToken wanted = %v  but got = %v ", tt.wantError, err)
				return
			}

			err = checkedSystem.IsTokenCorrect(token)
			if (err != nil) != tt.wantError && tt.errorPart == 3 {
				t.Errorf("TestTokens() IsTokenCorrect wanted = %v  but got = %v ", tt.wantError, err)
				return
			}

		})
	}
}
