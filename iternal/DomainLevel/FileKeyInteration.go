package DomainLevel

import "crypto/rsa"

type FileInfo interface {
	ConvertToBytesFileInfo(string, []byte) ([]byte, error)
	GetRealNameFile(string) string
	ProcessingFileParameters(string) (string, error)
	GenerateShortFileName() string
	EncryptData([]byte, *rsa.PublicKey) ([]byte, error)
	DecryptFileInfo([]byte, *rsa.PrivateKey, *rsa.PrivateKey) ([]byte, string, error)
}
