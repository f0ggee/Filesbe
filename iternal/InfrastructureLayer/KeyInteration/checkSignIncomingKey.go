package KeyInteration

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"log/slog"
)

type EncryptionKey struct{}

func (*EncryptionKey) DecryptAesKey(RsaKey []byte, aesKey []byte) ([]byte, error) {

	slog.Info("DecryptAesKey", "RsaKey", RsaKey)
	RsaKeyPrivate, err := x509.ParsePKCS1PrivateKey(RsaKey)
	if err != nil {
		slog.Error("Error Parsing RsaKey", "Func: DecryptAesKey", "Error", err)
		return nil, err
	}

	return rsa.DecryptOAEP(sha256.New(), rand.Reader, RsaKeyPrivate, aesKey, nil)
}

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
