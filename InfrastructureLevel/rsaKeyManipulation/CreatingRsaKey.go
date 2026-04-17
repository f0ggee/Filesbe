package rsaKeyManipulation

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"log/slog"
	"runtime/debug"
)

type RsaKeyManipulation struct{}

func (r *RsaKeyManipulation) GenerateRsaKey() []byte {
	slog.Info("Generating new RSA key")
	RsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		slog.Group("Error generating RSA key",
			slog.String("error", err.Error()),
			slog.String("stack", string(debug.Stack())),
			slog.Bool("Generating Rsa key ERROR", false))
		return nil
	}
	slog.Info("Generating Rsa key", "Success", true)
	return x509.MarshalPKCS1PrivateKey(RsaKey)

}
