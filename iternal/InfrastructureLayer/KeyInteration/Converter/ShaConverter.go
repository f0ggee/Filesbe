package Converter

import "crypto/sha256"

func (k KeyConverter) ConvertDataToHash(plainText []byte, aesKey []byte) []byte {
	shaHash := sha256.New()

	shaHash.Write(plainText)
	shaHash.Write(aesKey)
	return shaHash.Sum(nil)
}
