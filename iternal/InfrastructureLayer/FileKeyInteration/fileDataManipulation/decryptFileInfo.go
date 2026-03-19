package fileDataManipulation

import (
	"Kaban/iternal/Dto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
)

type FileDataManipulation struct{}

func (*FileDataManipulation) DecryptFileInfo(FileInfoIntoBytes []byte, key []byte, oldKey []byte) ([]byte, string, error) {

}
