package Handlers

import (
	"encoding/hex"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var ControlPrivateKeyStruct struct {
	MasterServerPublicKeyBytes []byte
	OurPrivateKeyIntoBytes     []byte
}
var Bucket string

func init() {

	err := godotenv.Load("iternal/Service/.env")
	if err != nil {
		slog.Error("cannot load env file", err.Error())
		return

	}
	Bucket = os.Getenv("BUCKET")
	PublickKeyIntoBytes, err := hex.DecodeString(os.Getenv("PublickKeyMasterServer"))
	if err != nil {
		slog.Error("Error decode publickKeyMasterServer", "Error", err.Error())
		return
	}
	OurPrivateKeyIntoBytes, err := hex.DecodeString(os.Getenv("Server2SecretKey"))
	if err != nil {
		slog.Error("Error decode Server1SecretKey", "Error", err.Error())
		return
	}
	ControlPrivateKeyStruct.MasterServerPublicKeyBytes = PublickKeyIntoBytes
	ControlPrivateKeyStruct.OurPrivateKeyIntoBytes = OurPrivateKeyIntoBytes

}
