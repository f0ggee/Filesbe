package Decription

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

type DecryptionData struct{}
