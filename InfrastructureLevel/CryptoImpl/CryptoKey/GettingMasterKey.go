package CryptoKey

import (
	"encoding/hex"
	"log/slog"
)

func (c CryptoManging) GetMasterKey() []byte {
	BytesMasterServerPrivateKey, err := hex.DecodeString(c.M.MasterServerSecret)
	if err != nil {
		slog.Error("Error decoding master server private key", "Error", err.Error())
		return nil
	}
	return BytesMasterServerPrivateKey
}
