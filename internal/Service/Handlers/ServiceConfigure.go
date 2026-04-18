package Handlers

import (
	"encoding/hex"
	"log/slog"
	"os"
)

type ControlPrivateKeyStruct struct {
	MasterServerPublicKeyBytes []byte
	OurPrivateKeyIntoBytes     []byte
}

var Bucket string

func ConfigureKeyData() {
	s := *new(ControlPrivateKeyStruct)
	Bucket = os.Getenv("BUCKET")
	PublickKeyIntoBytes, err := hex.DecodeString(os.Getenv("Public_Key_Master_Server"))
	if err != nil {
		slog.Error("Error decode Public_Key_Master_Server", "Error", err.Error())
		return
	}
	OurPrivateKeyIntoBytes, err := hex.DecodeString(os.Getenv("Our_Private_Key"))
	if err != nil {
		slog.Error("Error decode Server1SecretKey", "Error", err.Error())
		return
	}

	s.MasterServerPublicKeyBytes = PublickKeyIntoBytes
	s.OurPrivateKeyIntoBytes = OurPrivateKeyIntoBytes

}
