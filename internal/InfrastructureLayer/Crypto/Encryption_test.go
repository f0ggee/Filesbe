package Crypto_test

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/Crypto"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io"
	"testing"
)

var InputDataReader = bytes.NewReader([]byte("kP7#xQ2mZ9vLtR4wB8nY1cF6hJ3sD0gA5eU@A" +
	"dsdadasdn2"))
var InputData = []byte("kP7#xQ2mZ9vLtR4wB8nY1cF6hJ3sD0gA5eU@")

func TestCrypto(t *testing.T) {

	type CryptoData struct {
		Key   []byte
		nonce []byte
	}
	Input := []struct {
		name       string
		testData   []byte
		Encryption DomainLevel.CryptoMaker
		D          CryptoData
	}{
		{
			name:       "test1",
			testData:   InputData,
			Encryption: Crypto.NewAesCtr(),
			D: CryptoData{
				Key: func() []byte {
					key := make([]byte, 32)
					if _, err := io.ReadFull(rand.Reader, key); err != nil {
						panic(err)
					}
					return key
				}(),
				nonce: func() []byte {
					key := make([]byte, Crypto.NewAesCtr().GetRequiredOverheadSize())
					if _, err := io.ReadFull(rand.Reader, key); err != nil {
						panic(err)
					}
					return key
				}(),
			},
		},
		{
			name:       "test2",
			testData:   InputData,
			Encryption: Crypto.NewRsaEncryption(),
			D: CryptoData{
				Key: func() []byte {
					key, _ := rsa.GenerateKey(rand.Reader, 2048)
					return x509.MarshalPKCS1PrivateKey(key)
				}(),
				nonce: func() []byte {
					key := make([]byte, Crypto.NewRsaEncryption().GetRequiredOverheadSize())
					if _, err := io.ReadFull(rand.Reader, key); err != nil {
						panic(err)
					}
					return key
				}(),
			},
		},
		{
			name:       "test3",
			testData:   InputData,
			Encryption: Crypto.NewAesGcmEncryption(),
			D: CryptoData{
				Key: func() []byte {
					key := make([]byte, 32)
					if _, err := io.ReadFull(rand.Reader, key); err != nil {
						panic(err)
					}
					return key
				}(),
				nonce: func() []byte {
					key := make([]byte, Crypto.NewAesGcmEncryption().GetRequiredOverheadSize())
					if _, err := io.ReadFull(rand.Reader, key); err != nil {
						panic(err)
					}
					return key
				}(),
			},
		},
	}

	for _, tt := range Input {
		t.Run(tt.name, func(t *testing.T) {
			encrypt, err := tt.Encryption.MakeCrypto(tt.D.Key, tt.D.nonce)
			if err != nil {
				t.Errorf("Error crypto %v", err)
			}
			cipherText, err := encrypt.Encrypt(tt.testData)
			if err != nil {
				t.Errorf("Error crypto %v", err)

			}
			decrypt, err := tt.Encryption.MakeCrypto(tt.D.Key, tt.D.nonce)
			if err != nil {
				t.Errorf("Error crypto %v", err)
			}
			text, err := decrypt.Decrypt(cipherText)
			if err != nil {
				t.Errorf("Error crypto %v", err)
			}
			if !bytes.Equal(text, tt.testData) {
				t.Errorf("The input data %v and a decrypted text %v aren't alike", tt.testData, text)
			}
		})
	}

}
