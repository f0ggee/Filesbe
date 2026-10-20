package Application

import "errors"

var ErrorDecryptFile = errors.New("error to decrypt data")
var ErrorFileSizeBig = errors.New("the File's size is bigger than the default size")

const (
	FileMaxSize = 500000000
)
