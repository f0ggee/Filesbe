package send

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Drivers/Crypto"
	"Kaban/internal/Drivers/Parsers"
	"Kaban/internal/Drivers/RepoEncrypterKeys"
	"Kaban/internal/Protocol"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"testing"
	"time"
)

func TestFirstExchange_GetEncryptedPacket(t *testing.T) {
	type fields struct {
		CryptoGenerate DomainLevel.CryptoGenerating
		Encoder        DomainLevel.Encoder
		Aes            DomainLevel.CryptoMaker
		Rsa            DomainLevel.CryptoMaker
		Sign           DomainLevel.SignatureMaker
	}
	defaultFields := fields{
		CryptoGenerate: Crypto.GetNewGenerating(),
		Encoder:        Parsers.GetNewParsing(),
		Aes:            Crypto.NewAesGcmEncryption(),
		Rsa:            Crypto.NewRsaEncryption(),
		Sign:           Crypto.NewRsaSigner(),
	}
	tt := struct {
		name    string
		fields  fields
		wantErr bool
	}{

		name:    "test1",
		fields:  defaultFields,
		wantErr: false,
	}

	masterKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	serverKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	DomainLevel.GetCryptoKeys = RepoEncrypterKeys.NewTestPrepareData(x509.MarshalPKCS1PublicKey(&masterKey.PublicKey), x509.MarshalPKCS1PrivateKey(serverKey))
	f := FirstExchange{
		CryptoGenerate: tt.fields.CryptoGenerate,
		Encoder:        tt.fields.Encoder,
		Aes:            tt.fields.Aes,
		Rsa:            tt.fields.Rsa,
		Sign:           tt.fields.Sign,
	}
	got, err := f.GetEncryptedPacket()
	if (err != nil) != tt.wantErr {
		t.Errorf("GetEncryptedPacket() error = %v, wantErr %v", err, tt.wantErr)
		return
	}

	err = check(got, masterKey)
	if (err != nil) != tt.wantErr {
		t.Errorf("GetEncryptedPacket() error = %v, wantErr %v", err, tt.wantErr)
		return
	}
}

func check(Packet []byte, masterKey *rsa.PrivateKey) error {
	var IncomePacket Protocol.OutPacket

	f := struct {
		Decoder Parsers.Parsing
		Rsa     DomainLevel.CryptoMaker
		Aes     DomainLevel.CryptoMaker
	}{
		Decoder: Parsers.GetNewParsing(),
		Rsa:     Crypto.NewRsaEncryption(),
		Aes:     Crypto.NewAesGcmEncryption(),
	}
	err := f.Decoder.Decode(&IncomePacket, Packet)
	if err != nil {
		return err
	}

	rsaDecrypter, err := f.Rsa.MakeCrypto(x509.MarshalPKCS1PrivateKey(masterKey), nil)
	if err != nil {
		return err
	}

	decryptedAesKey, err := rsaDecrypter.Decrypt(IncomePacket.AesKeyData)
	if err != nil {
		return err
	}

	aesCrypto, err := f.Aes.MakeCrypto(decryptedAesKey, nil)
	if err != nil {
		return err
	}
	data, err := aesCrypto.Decrypt(IncomePacket.CipherData)
	if err != nil {
		return err
	}

	GrpcStruct := Protocol.MasterServerEnvelope{
		Time:             time.Time{},
		ServerName:       nil,
		SignedServerName: nil,
	}
	err = Parsers.GetNewParsing().Decode(&GrpcStruct, data)
	if err != nil {
		return err
	}
	signature, err := Crypto.NewRsaSigner().MakeSigner(DomainLevel.GetCryptoKeys.GetServerKey())
	if err != nil {
		return err
	}
	err = signature.CheckSign(GrpcStruct.SignedServerName, GrpcStruct.ServerName)
	if err != nil {
		return err
	}

	return nil

}
