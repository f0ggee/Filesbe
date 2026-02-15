package encryptionKeyInteration

import "crypto/sha256"

type EncryptionKey struct{}

func (*EncryptionKey) ConvertDataToHash(plainText []byte, aesKey []byte) []byte {
	shaHash := sha256.New()

	shaHash.Write(plainText)
	shaHash.Write(aesKey)
	return shaHash.Sum(nil)
}
