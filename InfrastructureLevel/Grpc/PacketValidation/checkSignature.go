package PacketValidation

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"log/slog"
)

func (s *ValidatePacketData) CheckSignature(Sign []byte, PublicServer []byte, ServerName []byte) error {
	PublicKey, err := x509.ParsePKCS1PublicKey(PublicServer)
	if err != nil {
		slog.Error("Error parsing public key", "Error", err.Error())
		return err
	}
	Hash := sha256.Sum256(ServerName)
	err = rsa.VerifyPKCS1v15(PublicKey, crypto.SHA256, Hash[:], Sign)
	if err != nil {
		slog.Error("Error validating the signature of the calling function ", "Error", err.Error())
		return err
	}
	return nil
}
