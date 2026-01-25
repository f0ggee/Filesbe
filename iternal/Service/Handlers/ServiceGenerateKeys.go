package Handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"log/slog"
	"sync"
)

var NewPrivateKey *rsa.PrivateKey
var OldPrivateKey *rsa.PrivateKey

var Mut sync.RWMutex

//SwapKeys generates a  pair keys

func SwapKeys() {
	Mut.Lock()

	OldPrivateKey = NewPrivateKey
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		slog.Error("Error generate key", err)
		return
	}

	NewPrivateKey = privateKey

	Mut.Unlock()
}
