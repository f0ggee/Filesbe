package Crypto

import (
	"Kaban/internal/DomainLevel"
	"crypto/rand"
	"fmt"
	"io"
	"testing"
)

func TestAesCtr_MakeCrypto(t *testing.T) {

	type args struct {
		key []byte
		i   []byte
	}
	key := []byte("examplekey123456")
	// 16 bytes key for AES-128.
	x := NewAesCtr()
	allData := make([]byte, x.GetRequiredRandomSize())
	if _, err := io.ReadFull(rand.Reader, allData); err != nil {
		panic(err)
	}

	tests := []struct {
		name      string
		inputData []byte
		args      args
		want      DomainLevel.Crypto
		wantErr   bool
	}{
		{
			name: "test1",
			args: args{
				key: key,
				i:   allData,
			},
			inputData: []byte("Test1"),
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type InputsData struct {
				Dats []byte
			}

			datas := make([]InputsData, 10000)
			encrypter, _ := NewAesCtr().MakeCrypto(tt.args.key, tt.args.i)
			for i := 0; i < 10000; i++ {
				encryptedData, _ := encrypter.Encrypt([]byte(rand.Text()))
				datas[i].Dats = encryptedData
			}

			decrypter, _ := NewAesCtr().MakeCrypto(tt.args.key, tt.args.i)
			for i := range datas {

				newData := make([]byte, len(datas[i].Dats))
				copy(newData, datas[i].Dats)
				decryptedData, _ := decrypter.Decrypt(newData)
				fmt.Println(string(decryptedData))

			}
			//encrypter, _ := NewAesCtr().MakeCrypto(tt.args.key, tt.args.i)
			//encryptedData, _ := encrypter.Encrypt([]byte("Hello!"))
			//decrypter, _ := NewAesCtr().MakeCrypto(tt.args.key, tt.args.i)
			//decryptedData, _ := decrypter.Decrypt(encryptedData)
			//fmt.Println(string(decryptedData))
			//
			////if !bytes.Equal(tt.inputData, decryptedData) {
			////	t.Errorf("Invalid data got =%v  wanted = %v", decryptedData, tt.inputData)
			////}

		})
	}
}
