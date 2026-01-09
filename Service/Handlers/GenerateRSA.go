package Handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"log/slog"
)

var NewPrivateKey *rsa.PrivateKey
var OldPrivateKey *rsa.PrivateKey

func SwapKeys() {

	OldPrivateKey = NewPrivateKey
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		slog.Error("Error generate key", err)
		return
	}

	NewPrivateKey = privateKey

}
