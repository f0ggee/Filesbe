package PacketChecking

import (
	"Kaban/iternal/DomainLevel"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"log/slog"
)

type Validating struct {
	Decrypter DomainLevel.KeyInteraction
}

func (v *Validating) CheckSign(sign []byte, Key []byte, Public []byte) error {

	PublickKey, err := x509.ParsePKCS1PublicKey(Public)
	if err != nil {
		slog.Error("Error Parsing Public Key", "Func: CheckSign", "Error", err.Error())
		return err
	}

	hash := sha256.Sum256(Key)

	err = rsa.VerifyPKCS1v15(PublickKey, crypto.SHA256, hash[:], sign)
	if err != nil {
		slog.Error("Error Verifying  ", "Func: CheckSign", "Error", err.Error())
		return err
	}
	return nil

}
