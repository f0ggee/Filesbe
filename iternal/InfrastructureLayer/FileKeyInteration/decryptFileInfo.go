package FileKeyInteration

import (
	"Kaban/iternal/Dto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log/slog"
)

func (*FileInfoController) DecryptFileInfo(FileInfoIntoBytes []byte, key *rsa.PrivateKey, oldKey *rsa.PrivateKey) ([]byte, string, error) {

	decryptFileInfo, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, FileInfoIntoBytes, nil)

	//switch {
	//case strings.Contains(fmt.Sprint(err), "decryption error"):
	//	slog.Error("Key is old")
	//	decryptFileInfo, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, oldKey, FileInfoIntoBytes, nil)
	//	if err != nil {
	//		slog.Error("Error also decrypt with an old key ", err)
	//		return nil, "", err
	//	}
	//
	//}
	if err != nil {
		slog.Error("Error in decrypt file info", "Error", err)
		return nil, "", err
	}

	sa := &Dto.FileDescription{
		FileName: "",
		AesKey:   "",
	}
	err = json.Unmarshal(decryptFileInfo, &sa)
	if err != nil {
		slog.Error("Error unmarshal aes", "ERR", err)
		return nil, "", err
	}

	aesKeyIntoByte, err := hex.DecodeString(sa.AesKey)
	if err != nil {
		slog.Error("Error decode aes key into string", err)
		return nil, "", err
	}

	return aesKeyIntoByte, sa.FileName, nil
}
