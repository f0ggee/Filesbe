package Crypto

import (
	"Kaban/internal/DomainLevel"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"testing"
)

func TestSigners(t *testing.T) {

	var inputData = []byte("Test")

	tests := []struct {
		name      string
		key       []byte
		Method    DomainLevel.SignatureMaker
		wantError bool
	}{
		{
			name: "test1",
			key: func() []byte {
				rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
				if err != nil {
					panic(err)
				}
				return x509.MarshalPKCS1PrivateKey(rsaKey)

			}(),
			Method:    NewRsaSigner(),
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signer, err := tt.Method.MakeSigner(tt.key)
			if (err != nil) != tt.wantError {
				t.Errorf("TestSigners() wanted = %v  but got = %v ", tt.wantError, err)
				return
			}
			sign, err := signer.MakeSign(inputData)
			if (err != nil) != tt.wantError {
				t.Errorf("TestSigners() wanted = %v  but got = %v ", tt.wantError, err)
				return
			}
			err = signer.CheckSign(sign, inputData)
			if (err != nil) != tt.wantError {
				t.Errorf("TestSigners() wanted = %v  but got = %v ", tt.wantError, err)
				return
			}
		})
	}
}
