package Application

import (
	"errors"
)

var ErrorDecryptFile = errors.New("error to decrypt data")
var ErrorFileSizeBig = errors.New("the File's size is bigger than the default size")
var ErrorWrongCryptoData = errors.New("crypto data is wrong")
var ErrorUploadingEncrypt = errors.New("an unexpected error occurred while loading encryption")

const (
	FileMaxSize = 500 * 1024 * 1024
	ChunkSize   = 30 * 1024 * 1024
)
