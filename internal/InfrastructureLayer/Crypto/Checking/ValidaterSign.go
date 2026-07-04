package Checking

import (
	"Kaban/internal/DomainLevel"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"log/slog"
)

func (c *Validating) CheckSignKey(data DomainLevel.CheckSignKeyIncomingData) error {
	publicKeyMasterServer, err := x509.ParsePKCS1PublicKey(data.MasterPublicKey)
	if err != nil {
		slog.Error("CheckSignKey; Error marshalling public key", "ERROR", err.Error())
		return err
	}

	err = rsa.VerifyPKCS1v15(publicKeyMasterServer, crypto.SHA256, data.Hash, data.Sign)
	if err != nil {
		slog.Error("CheckSignKey; Error verifying signature", "Error", err.Error())
		return err
	}
	return nil
}
