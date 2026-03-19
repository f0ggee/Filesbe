package Encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

type Encrypter struct {
}
