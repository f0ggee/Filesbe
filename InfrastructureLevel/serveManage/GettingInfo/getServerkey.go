package GettingInfo

import (
	"encoding/hex"
	"log/slog"
	"os"
	"time"
)

var MappingHash = make(map[[32]byte]time.Time)

func (s *SeverManage) GetServerKey(Num int) []byte {
	switch Num {
	case 1:

		rsaKey, err := hex.DecodeString(os.Getenv("SERVER_1"))
		if err != nil {
			slog.Error("Error in getting the server key", "ServerNumber", Num)
			return nil
		}

		return rsaKey

	case 2:
		rsaKey, err := hex.DecodeString(os.Getenv("SERVER_2"))
		if err != nil {
			slog.Error("Error in getting the server key", "ServerNumber", Num)
			return nil
		}
		return rsaKey

	}

	slog.Info("Couldn't find the server key", "ServerNumber", Num)
	return nil
}
