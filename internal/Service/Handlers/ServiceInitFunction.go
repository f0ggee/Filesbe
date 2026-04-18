package Handlers

import (
	"encoding/hex"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type ControlPrivateKeyStruct struct {
	MasterServerPublicKeyBytes []byte
	OurPrivateKeyIntoBytes     []byte
}

var Bucket string

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("cannot load env file", "Error", err.Error())

	}
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
