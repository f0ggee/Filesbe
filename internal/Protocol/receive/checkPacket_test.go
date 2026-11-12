package receive

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Drivers/Crypto"
	"Kaban/internal/Drivers/Parsers"
	"Kaban/internal/Drivers/RepoEncrypterKeys"
	"Kaban/internal/Protocol"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io"
	"testing"
	"time"
)

func TestGetNewKey_GetPacketData(t *testing.T) {
	type fields struct {
		Decoder  DomainLevel.Decoder
		Aes      DomainLevel.CryptoMaker
		Rsa      DomainLevel.CryptoMaker
		Validate DomainLevel.SignatureMaker
	}
	type args struct {
		packet []byte
	}
	defaulFields := fields{
		Decoder:  Parsers.GetNewParsing(),
		Aes:      Crypto.NewAesGcmEncryption(),
		Rsa:      Crypto.NewRsaEncryption(),
		Validate: Crypto.NewRsaSigner(),
	}

	defaultArgs := args{
		packet: masterServerTest(),
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  defaulFields,
			args:    defaultArgs,
			wantErr: false,
		},
		{
			name:   "test2",
			fields: defaulFields,
			args: args{
				packet: badEncryptionRsa(),
			},
			wantErr: true,
		},
		{
			name:   "test3",
			fields: defaulFields,
			args: args{
				packet: badEncoding(),
			},
			wantErr: true,
		},
		{
			name:   "test4",
			fields: defaulFields,
			args: args{
				packet: badTime(),
			},
			wantErr: true,
		},
		{
			name:   "test5",
			fields: defaulFields,
			args: args{
				packet: badSignRsa(),
			},
			wantErr: true,
		},
		{
			name:   "test6",
			fields: defaulFields,
			args: args{
				packet: badPayload(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := GetNewKey{
				Decoder:  tt.fields.Decoder,
				Aes:      tt.fields.Aes,
				Rsa:      tt.fields.Rsa,
				Validate: tt.fields.Validate,
			}
			_, err := f.GetPacketData(tt.args.packet)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPacketData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func badEncoding() []byte {
	p := Protocol.InEnvelope{
		Sign:    []byte("eqwewq"),
		RsaKey:  []byte("ehsad"),
		T1:      1000,
		TimeNow: time.Now(),
	}
	d, _ := Parsers.GetNewParsing().Encode(p)
	return d
}
func badPayload() []byte {
	masterKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	serverKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	playLoadBytes := []byte(rand.Text())

	sign, err := Crypto.NewRsaSigner().MakeSigner(x509.MarshalPKCS1PrivateKey(masterKey))
	if err != nil {
		panic(err)
	}
	signature, err := sign.MakeSign(playLoadBytes)
	if err != nil {
		panic(err)
	}
	data, _ := Parsers.GetNewParsing().Encode(Protocol.InEnvelope{
		Sign:    signature,
		RsaKey:  []byte("DatataDatataDatataDatataDatataDatata"),
		T1:      100,
		TimeNow: time.Now(),
	})

	aesKey := make([]byte, 32)
	io.ReadFull(rand.Reader, aesKey)
	_aesCrypto, _ := Crypto.NewAesGcmEncryption().MakeCrypto(aesKey, nil)
	encryptedPacket, _ := _aesCrypto.Encrypt(data)

	_rsaCrypto, _ := Crypto.NewRsaEncryption().MakeCrypto(x509.MarshalPKCS1PrivateKey(serverKey), nil)

	encryptedKey, err := _rsaCrypto.Encrypt(aesKey)
	if err != nil {
		panic(err)
	}

	DomainLevel.GetCryptoKeys = RepoEncrypterKeys.NewTestEncrypterKeys([]byte("e"))
	DomainLevel.ServerKeys = DomainLevel.NewServerKeys{
		MasterKeyPublic:  x509.MarshalPKCS1PublicKey(&masterKey.PublicKey),
		WorkerPrivateKey: x509.MarshalPKCS1PrivateKey(serverKey),
	}

	Packet := Protocol.OutPacket{
		AesKeyData: encryptedKey,
		CipherData: encryptedPacket,
	}

	packetBytes, _ := Parsers.GetNewParsing().Encode(Packet)

	return packetBytes
}
func badSignRsa() []byte {
	masterKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	serverKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	badKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	playLoadBytes := []byte(rand.Text())

	sign, err := Crypto.NewRsaSigner().MakeSigner(x509.MarshalPKCS1PrivateKey(badKey))
	if err != nil {
		panic(err)
	}
	signature, err := sign.MakeSign(playLoadBytes)
	if err != nil {
		panic(err)
	}
	data, _ := Parsers.GetNewParsing().Encode(Protocol.InEnvelope{
		Sign:    signature,
		RsaKey:  playLoadBytes,
		T1:      100,
		TimeNow: time.Now(),
	})

	aesKey := make([]byte, 32)
	io.ReadFull(rand.Reader, aesKey)
	_aesCrypto, _ := Crypto.NewAesGcmEncryption().MakeCrypto(aesKey, nil)
	encryptedPacket, _ := _aesCrypto.Encrypt(data)

	_rsaCrypto, _ := Crypto.NewRsaEncryption().MakeCrypto(x509.MarshalPKCS1PrivateKey(serverKey), nil)

	encryptedKey, err := _rsaCrypto.Encrypt(aesKey)
	if err != nil {
		panic(err)
	}

	DomainLevel.GetCryptoKeys = RepoEncrypterKeys.NewTestEncrypterKeys([]byte("e"))
	DomainLevel.ServerKeys = DomainLevel.NewServerKeys{
		MasterKeyPublic:  x509.MarshalPKCS1PublicKey(&masterKey.PublicKey),
		WorkerPrivateKey: x509.MarshalPKCS1PrivateKey(serverKey),
	}

	Packet := Protocol.OutPacket{
		AesKeyData: encryptedKey,
		CipherData: encryptedPacket,
	}

	packetBytes, _ := Parsers.GetNewParsing().Encode(Packet)

	return packetBytes
}
func badEncryptionRsa() []byte {
	badKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	data := []byte(rand.Text())
	_rsa, _ := Crypto.NewRsaEncryption().MakeCrypto(x509.MarshalPKCS1PrivateKey(badKey), nil)
	encryptedData, _ := _rsa.Encrypt(data)

	outcomePacket := Protocol.OutPacket{
		AesKeyData: encryptedData,
		CipherData: data,
	}
	packet, _ := Parsers.GetNewParsing().Encode(outcomePacket)
	return packet
}

func badTime() []byte {
	masterKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	serverKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	playLoadBytes := []byte(rand.Text())

	sign, err := Crypto.NewRsaSigner().MakeSigner(x509.MarshalPKCS1PrivateKey(masterKey))
	if err != nil {
		panic(err)
	}
	signature, err := sign.MakeSign(playLoadBytes)
	if err != nil {
		panic(err)
	}
	data, _ := Parsers.GetNewParsing().Encode(Protocol.InEnvelope{
		Sign:    signature,
		RsaKey:  playLoadBytes,
		T1:      100,
		TimeNow: time.Now().Add(-5 * time.Hour),
	})

	aesKey := make([]byte, 32)
	io.ReadFull(rand.Reader, aesKey)
	_aesCrypto, _ := Crypto.NewAesGcmEncryption().MakeCrypto(aesKey, nil)
	encryptedPacket, _ := _aesCrypto.Encrypt(data)

	_rsaCrypto, _ := Crypto.NewRsaEncryption().MakeCrypto(x509.MarshalPKCS1PrivateKey(serverKey), nil)

	encryptedKey, err := _rsaCrypto.Encrypt(aesKey)
	if err != nil {
		panic(err)
	}

	DomainLevel.GetCryptoKeys = RepoEncrypterKeys.NewTestEncrypterKeys([]byte("e"))
	DomainLevel.ServerKeys = DomainLevel.NewServerKeys{
		MasterKeyPublic:  x509.MarshalPKCS1PublicKey(&masterKey.PublicKey),
		WorkerPrivateKey: x509.MarshalPKCS1PrivateKey(serverKey),
	}

	Packet := Protocol.OutPacket{
		AesKeyData: encryptedKey,
		CipherData: encryptedPacket,
	}

	packetBytes, _ := Parsers.GetNewParsing().Encode(Packet)
	return packetBytes
}
func masterServerTest() []byte {
	masterKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	serverKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	playLoadBytes := []byte(rand.Text())

	sign, err := Crypto.NewRsaSigner().MakeSigner(x509.MarshalPKCS1PrivateKey(masterKey))
	if err != nil {
		panic(err)
	}
	signature, err := sign.MakeSign(playLoadBytes)
	if err != nil {
		panic(err)
	}
	data, _ := Parsers.GetNewParsing().Encode(Protocol.InEnvelope{
		Sign:    signature,
		RsaKey:  playLoadBytes,
		T1:      100,
		TimeNow: time.Now(),
	})

	aesKey := make([]byte, 32)
	io.ReadFull(rand.Reader, aesKey)
	_aesCrypto, _ := Crypto.NewAesGcmEncryption().MakeCrypto(aesKey, nil)
	encryptedPacket, _ := _aesCrypto.Encrypt(data)

	_rsaCrypto, _ := Crypto.NewRsaEncryption().MakeCrypto(x509.MarshalPKCS1PrivateKey(serverKey), nil)

	encryptedKey, err := _rsaCrypto.Encrypt(aesKey)
	if err != nil {
		panic(err)
	}

	DomainLevel.GetCryptoKeys = RepoEncrypterKeys.NewTestEncrypterKeys([]byte("e"))
	DomainLevel.ServerKeys = DomainLevel.NewServerKeys{
		MasterKeyPublic:  x509.MarshalPKCS1PublicKey(&masterKey.PublicKey),
		WorkerPrivateKey: x509.MarshalPKCS1PrivateKey(serverKey),
	}

	Packet := Protocol.OutPacket{
		AesKeyData: encryptedKey,
		CipherData: encryptedPacket,
	}

	packetBytes, _ := Parsers.GetNewParsing().Encode(Packet)

	return packetBytes
}
