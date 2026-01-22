package Handlers

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var Bucket string

func init() {

	err := godotenv.Load("iternal/Service/.env")
	if err != nil {
		slog.Error("cannot load env file", err)
		return

	}
	Bucket = os.Getenv("BUCKET")
}
