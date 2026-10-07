//This file contains functions that work with packets income/outcome.
//It encrypts a packer or decrypts it.

package send

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/Protocol"
	"crypto/rand"
	"log/slog"
	"os"
	"time"

	"github.com/awnumar/memguard"
)

type FirstExchange struct {
	CryptoGenerate DomainLevel.CryptoGenerating
	Encoder        DomainLevel.Encoder
	Aes            DomainLevel.CryptoMaker
	Rsa            DomainLevel.CryptoMaker
}

// GetEncryptedPacket prepares and encrypts a packet to send it to a Master server
func (f FirstExchange) GetEncryptedPacket() ([]byte, error) {
	var serverName = os.Getenv("server_name")
	signedServerName, err := f.CryptoGenerate.GenerateSignature([]byte(serverName), DomainLevel.ServerKeys.GerOurPrivateKey())
	if err != nil {
		return nil, err
	}
	var GrpcStruct = Dto.GrpcOutComingPacketDetails{
		Time:             time.Now(),
		ServerName:       []byte(serverName),
		SignedServerName: signedServerName,
	}

	convertedData, err := f.Encoder.Encode(&GrpcStruct)
	if err != nil {
		return nil, err
	}

	var ACD = f.Aes.GetRequiredOverheadSize()
	aesKey, err := memguard.NewBufferFromReader(rand.Reader, 32+ACD)
	if err != nil {
		slog.Error("GetEncryptedPacket: error to generate data", "ERROR", err)
		return nil, Protocol.ErrorPrepareData
	}
	defer aesKey.Destroy()

	aesMaker, err := f.Aes.MakeCrypto(aesKey.Data()[ACD:], aesKey.Data()[:ACD])
	if err != nil {
		return nil, err
	}
	encryptedPacket, err := aesMaker.Encrypt(convertedData)
	if err != nil {
		return nil, err
	}
	rsaMaker, err := f.Rsa.MakeCrypto(DomainLevel.ServerKeys.GetMasterPublicKey(), nil)
	if err != nil {
		return nil, err
	}
	encryptedFileInfo, err := rsaMaker.Encrypt(aesKey.Data()[ACD:])
	if err != nil {
		return nil, err
	}

	var EncryptedOutData = Dto.GrpcOutComingPacketForSending{
		AesKeyData: encryptedFileInfo,
		CipherData: encryptedPacket,
	}

	convertedCryptoPacket, err := f.Encoder.Encode(&EncryptedOutData)
	if err != nil {
		return nil, err
	}

	return convertedCryptoPacket, nil
}
