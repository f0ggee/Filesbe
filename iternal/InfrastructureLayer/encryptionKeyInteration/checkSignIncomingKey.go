package encryptionKeyInteration

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"log/slog"
)

func (*EncryptionKey) CheckSignIncomingKey(sign []byte, hash []byte, PublicKeyMasterServer []byte) error {

	publicKeyMasterServer, err := x509.ParsePKCS1PublicKey(PublicKeyMasterServer)
	if err != nil {
		slog.Error("Error marshalling public key", "Error", err.Error())
		return err
	}
	err = rsa.VerifyPKCS1v15(publicKeyMasterServer, crypto.SHA256, hash, sign)
	if err != nil {
		slog.Error("Error verifying signature", "Error", err.Error())
		return err
	}
	return nil

}
