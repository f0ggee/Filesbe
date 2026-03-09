package KeyInteration

import (
	"crypto/sha256"
	"encoding/json"
)

type EncryptionKey struct{}

func (*EncryptionKey) JsonConverter(data any) ([]byte, error) {

	JsonDataType, err := json.Marshal(&data)

	if err != nil {
		return nil, err
	}

	return JsonDataType, err
}

func (*EncryptionKey) ConvertDataToHash(plainText []byte, aesKey []byte) []byte {
	shaHash := sha256.New()

	shaHash.Write(plainText)
	shaHash.Write(aesKey)
	return shaHash.Sum(nil)
}
