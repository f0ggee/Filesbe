package Crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io"
	"testing"
)

func TestAesCtr_MakeCrypto1(t *testing.T) {

	var NonceGeneration func() []byte

	NonceGeneration = func() []byte {
		nonce := make([]byte, aes.BlockSize)
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			panic(err)
		}
		return nonce
	}
	key := []byte("examplekey123456") // 16 bytes key1 for AES-128.
	type args struct {
		key []byte
		i   []byte
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		inputData []byte
	}{
		{
			name: "test1",
			args: args{
				key: key,
				i:   NonceGeneration(),
			},
			wantErr:   false,
			inputData: []byte(rand.Text()),
		},
		{
			name: "test2",
			args: args{
				key: nil,
				i:   NonceGeneration(),
			},
			wantErr:   true,
			inputData: []byte(rand.Text()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAesCtr().MakeCrypto(tt.args.key, tt.args.i)
			if err != nil && !tt.wantErr {
				t.Errorf("MakeCrypto() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			encryptedData, err := got.Encrypt(tt.inputData)
			if err != nil && !tt.wantErr {
				t.Errorf("MakeCrypto() error = %v, wantErr %v", err, tt.wantErr)

			}

			decryptedGot, err := NewAesCtr().MakeCrypto(tt.args.key, tt.args.i)
			if err != nil && !tt.wantErr {
				t.Errorf("MakeCrypto() error = %v, wantErr %v", err, tt.wantErr)

			}
			decryptedData, err := decryptedGot.Decrypt(encryptedData)
			if err != nil && !tt.wantErr {
				t.Errorf("MakeCrypto() error = %v, wantErr %v", err, tt.wantErr)

			}
			if !bytes.Equal(tt.inputData, decryptedData) {
				t.Errorf("The input data = %v and decrypted data = %v aren't alike", string(tt.inputData), string(decryptedData))
				return
			}

		})
	}
}

func TestAesEncryption_MakeCrypto(t *testing.T) {
	var NonceGeneration func() []byte

	NonceGeneration = func() []byte {
		nonce := make([]byte, NewAesEncryption().GetRequiredOverheadSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			panic(err)
		}
		return nonce
	}
	NoncegenerationBad := func() []byte {
		nonce := make([]byte, 16)
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			panic(err)
		}
		return nonce
	}
	key := []byte("examplekey123456")
	type args struct {
		bytes []byte
		i     []byte
	}
	tests := []struct {
		name string

		args      args
		wantErr   bool
		InputData []byte
	}{
		{
			name: "test1",
			args: args{
				bytes: key,
				i:     NonceGeneration(),
			},
			wantErr:   false,
			InputData: []byte(rand.Text()),
		},
		{
			name: "test2",
			args: args{
				bytes: key,
				i:     NoncegenerationBad(),
			},
			wantErr:   true,
			InputData: []byte(rand.Text()),
		},

		{
			name: "test3",
			args: args{
				bytes: nil,
				i:     nil,
			},
			wantErr:   true,
			InputData: nil,
		},
		{
			name: "test4",
			args: args{
				bytes: key,
				i:     NonceGeneration(),
			},
			wantErr:   false,
			InputData: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAesEncryption().MakeCrypto(tt.args.bytes, tt.args.i)
			if err != nil && !tt.wantErr {
				t.Errorf("MakeCrypto() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil && tt.wantErr {
				t.SkipNow()
			}

			encryptedData, err := got.Encrypt(tt.InputData)
			if err != nil && !tt.wantErr {
				t.Errorf("MakeCrypto() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil && tt.wantErr {
				t.SkipNow()
			}
			newGot2, err := NewAesEncryption().MakeCrypto(tt.args.bytes, tt.args.i)
			if err != nil && !tt.wantErr {
				t.Errorf("MakeCrypto() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil && tt.wantErr {
				t.SkipNow()
			}

			decryptedData, err := newGot2.Decrypt(encryptedData)
			if err != nil && !tt.wantErr {
				t.Errorf("MakeCrypto() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil && tt.wantErr {
				t.SkipNow()
			}
			if !bytes.Equal(decryptedData, tt.InputData) {
				t.Errorf("Decrypting want = %v got = %v", decryptedData, tt.InputData)
				return
			}
		})
	}
}

func TestRsaEncryption_MakeCrypto(t *testing.T) {

	generateKets := func() []byte {
		keys, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			panic(err)
		}
		return x509.MarshalPKCS1PrivateKey(keys)
	}
	keys := generateKets()
	pub := func(keys []byte) []byte {

		x, err := x509.ParsePKCS1PrivateKey(keys)
		if err != nil {
			panic(err)
		}
		return x509.MarshalPKCS1PublicKey(&x.PublicKey)
	}

	type args struct {
		key1        []byte
		determiner1 []byte
		key2        []byte
		determiner2 []byte
	}
	tests := []struct {
		name string
		args
		inputData []byte
		wantErr   bool
	}{
		{
			name: "test1",
			args: args{
				key1:        pub(keys),
				determiner1: nil,
				key2:        (keys),
				determiner2: nil,
			},
			inputData: []byte(rand.Text()),
			wantErr:   false,
		},
		{
			name:      "test2",
			args:      args{},
			inputData: nil,
			wantErr:   true,
		},
		{
			name: "test3",
			args: args{
				key1:        keys,
				determiner1: nil,
				key2:        nil,
				determiner2: nil,
			},
			inputData: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRsaEncryption().MakeCrypto(tt.args.key1, tt.args.determiner1)
			if err != nil && !tt.wantErr {
				t.Errorf("MakeCrypto() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr {
				t.Logf("MakeCrypto() we wantedErr = %v, got = %v", tt.wantErr, err)
				return
			}
			encryptedData, err := got.Encrypt(tt.inputData)
			if err != nil && !tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr {
				t.Logf("Encrypt() we wantedErr = %v, got = %v", tt.wantErr, err)
				return

			}
			newGot2, err := NewRsaEncryption().MakeCrypto(tt.args.key2, tt.args.determiner2)
			if err != nil && !tt.wantErr {
				t.Errorf("MakeCrypto() error = %v, wantErr %v", err, tt.wantErr)
				return

			}
			if err != nil && tt.wantErr {
				t.Errorf("MakeCrypto() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			decryptedData, err := newGot2.Decrypt(encryptedData)
			if err != nil && !tt.wantErr {
				t.Errorf("MakeCrypto() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !bytes.Equal(tt.inputData, decryptedData) {
				t.Errorf("Decrypting want = %v got = %v", decryptedData, tt.inputData)
				return
			}
			t.Logf(
				"Decrypter: input data = %v decrypted data = %v",
				string(tt.inputData),
				string(decryptedData),
			)

		})
	}
}
