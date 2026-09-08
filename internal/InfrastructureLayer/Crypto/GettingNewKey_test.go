package Crypto

import (
	"Kaban/internal/DomainLevel"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"testing"
)

func CheckSingTest(h Checking, Databyte, Hash, Key []byte) error {

	data := DomainLevel.CheckSignKeyIncomingData{
		Sign:            Databyte,
		Hash:            Hash,
		MasterPublicKey: Key,
	}
	err := h.CheckSign(data)
	if err != nil {
		return err
	}
	return nil
}
func DecryptAesKeyTest(sa *Decryption, key, data []byte) error {

	_, err := sa.DecryptAesKey(key, data)
	if err != nil {
		return err
	}
	return nil
}
func MakeRsaKey() []byte {

	rsa, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	return x509.MarshalPKCS1PrivateKey(rsa)
}

func BadEncrypt(sa Encrypter) []byte {

	rsae, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	rea, err := sa.EncryptFileInfo([]byte("It's test!"), &rsae.PublicKey)

	if err != nil {
		panic(err)
	}
	return rea
}
func GooEncrypt(sa Encrypter, e []byte) []byte {

	ez, err := x509.ParsePKCS1PrivateKey(e)
	if err != nil {
		panic(err)
	}

	sax, err := sa.EncryptFileInfo([]byte("It's test!"), &ez.PublicKey)
	if err != nil {
		panic(err)

	}
	return sax
}

func TestEe2(t *testing.T) {
	type args struct {
		Sa   *Decryption
		key  []byte
		data []byte
	}

	Realization := &Decryption{}
	encrypter := Encrypter{}
	OurKey := MakeRsaKey()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test_1",
			args:    args{Sa: Realization, key: OurKey, data: BadEncrypt(encrypter)},
			wantErr: true,
		},
		{
			name:    "Test_2",
			args:    args{Sa: Realization, key: OurKey, data: GooEncrypt(encrypter, OurKey)},
			wantErr: false,
		},
		{
			name: "Test_3",
			args: args{
				Sa:   Realization,
				key:  nil,
				data: nil,
			},

			wantErr: true,
		},
		{
			name: "Test_4",
			args: args{
				Sa:   Realization,
				key:  nil,
				data: GooEncrypt(encrypter, OurKey),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DecryptAesKeyTest(tt.args.Sa, tt.args.key, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("DecryptAesKeyTest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

const TESTDATA = "IT's test"

func GoodSign(key []byte) []byte {

	Key, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		panic(err)
	}

	sha := sha256.New()

	sha.Write([]byte(TESTDATA))
	Sing, err := rsa.SignPKCS1v15(rand.Reader, Key, crypto.SHA256, sha.Sum([]byte(nil)))
	if err != nil {
		panic("ESDCXz" + err.Error())
	}

	return Sing
}

func BadSing() []byte {

	Key, err := x509.ParsePKCS1PrivateKey(MakeRsaKey())
	if err != nil {
		panic("ee" + err.Error())
	}

	sha := sha256.New()
	sha.Write([]byte(TESTDATA))
	Data, err := rsa.SignPKCS1v15(rand.Reader, Key, crypto.SHA256, sha.Sum([]byte(nil)))
	if err != nil {
		panic("eesadwqw" + err.Error())
	}
	return Data
}
func Test_TestSign(t *testing.T) {
	type args struct {
		h    Checking
		Data []byte
		Hash []byte
		Key  []byte
	}

	OurKEy := MakeRsaKey()

	e, err := x509.ParsePKCS1PrivateKey(OurKEy)
	if err != nil {
		panic(err)
	}
	zs := x509.MarshalPKCS1PublicKey(&e.PublicKey)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test_1",
			args: args{
				h:    Checking{},
				Data: GoodSign(OurKEy),
				Hash: []byte(TESTDATA),
				Key:  zs,
			},
			wantErr: false,
		},
		{
			name: "Test_2",
			args: args{
				h:    Checking{},
				Data: GoodSign(OurKEy),
				Hash: []byte(TESTDATA),
				Key:  nil,
			},
			wantErr: true,
		},
		{
			name: "Test_3",
			args: args{
				h:    Checking{},
				Data: BadSing(),
				Hash: []byte(TESTDATA),
				Key:  zs,
			},
			wantErr: true,
		},
		{
			name: "Test_4",
			args: args{
				h:    Checking{},
				Data: GoodSign(OurKEy),
				Hash: []byte("CRYPTO"),
				Key:  zs,
			},
			wantErr: true,
		},
		{
			name: "Test_5",
			args: args{
				h:    Checking{},
				Data: BadSing(),
				Hash: nil,
				Key:  nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			SASHA := sha256.New()
			SASHA.Write(tt.args.Hash)
			if err := CheckSingTest(tt.args.h, tt.args.Data, SASHA.Sum([]byte(nil)), tt.args.Key); (err != nil) != tt.wantErr {
				t.Errorf("eee() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
