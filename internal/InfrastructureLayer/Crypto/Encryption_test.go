package Crypto

import (
	"crypto/aes"
	"crypto/rand"
	"fmt"
	"io"
	"testing"
)

func TestAesCtr_Encrypter(t *testing.T) {

	InputData := []byte("Hello!")
	key := []byte("examplekey123456") // 16 bytes key for AES-128.
	allData := make([]byte, aes.BlockSize, aes.BlockSize+len(key))
	if _, err := io.ReadFull(rand.Reader, allData); err != nil {
		panic(err)
	}
	allData = append(allData, key...)
	encrypter, err := NewAesCtr().MakeCrypto(allData, nil)
	if err != nil {
		panic(err)
	}

	EncryptedDataOuput, err := encrypter.Encrypt(InputData)
	if err != nil {
		panic(err)
	}

	alldataDecrypted := []byte{}
	alldataDecrypted = append(alldataDecrypted, EncryptedDataOuput[:aes.BlockSize]...)
	alldataDecrypted = append(alldataDecrypted, key...)
	decrypter, err := NewAesCtr().MakeCrypto(alldataDecrypted, nil)
	if err != nil {
		panic(err)
	}

	DecryptedData, err := decrypter.Decrypt(EncryptedDataOuput[aes.BlockSize:])
	if err != nil {
		panic(err)
	}
	fmt.Println(string(DecryptedData))
}
