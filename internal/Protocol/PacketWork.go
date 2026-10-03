//This file contains functions that work with packets income/outcome.
//It encrypts a packer or decrypts it.

package Protocol

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"crypto/rand"
	"crypto/sha256"
	"log/slog"
	"os"
	"time"

	"github.com/awnumar/memguard"
)

type FirstExchange struct {
	CryptoGenerate DomainLevel.CryptoGenerating
	Keys           DomainLevel.NewServerKeys
	Encoder        DomainLevel.Encoder
	Encrypter1     DomainLevel.CryptoMaker
	Encrypter2     DomainLevel.CryptoMaker
}

// GetEncryptedPacket prepares and encrypts a packet to send it to a Master server
func (f FirstExchange) GetEncryptedPacket() ([]byte, error) {
	var serverName = os.Getenv("server_name")
	signedServerName, err := f.CryptoGenerate.GenerateSignature([]byte(serverName), f.Keys.GerOurPrivateKey())
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

	//Additional crypto data
	var ACD = f.Encrypter1.GetRequiredOverheadSize()
	aesKey, err := memguard.NewBufferFromReader(rand.Reader, 32+ACD)
	if err != nil {
		slog.Error("GetEncryptedPacket: error to generate data", "ERROR", err)
		return nil, ErrorPrepareData
	}
	defer aesKey.Destroy()

	aesMaker, err := f.Encrypter1.MakeCrypto(aesKey.Data()[ACD:], aesKey.Data()[:ACD])
	if err != nil {
		return nil, err
	}
	encryptedPacket, err := aesMaker.Encrypt(convertedData)
	if err != nil {
		return nil, err
	}
	rsaMaker, err := f.Encrypter2.MakeCrypto(f.Keys.GetMasterPublicKey(), nil)
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

type GetNewKey struct {
	Keys           DomainLevel.NewServerKeys
	Decoder        DomainLevel.Decoder
	aes            DomainLevel.CryptoMaker
	rsa            DomainLevel.CryptoMaker
	CryptoValidate DomainLevel.CryptoValidating
	TempKey        RepoEncrypterKeys.Keys
}

// GetPacketData decrypts and cheks a packer and then returns data.
func (f GetNewKey) GetPacketData(packet []byte) (time.Duration, error) {
	var IncomePacket Dto.GrpcOutComingPacketForSending
	err := f.Decoder.Decode(&IncomePacket, packet)
	if err != nil {
		return 0, err
	}

	var ACD = f.aes.GetRequiredOverheadSize()

	rsaDecrypter, err := f.rsa.MakeCrypto(f.Keys.GerOurPrivateKey(), nil)
	if err != nil {
		return 0, err
	}

	decryptedAesKey, err := rsaDecrypter.Decrypt(IncomePacket.AesKeyData)
	if err != nil {
		return 0, err
	}

	aesCrypto, err := f.aes.MakeCrypto(decryptedAesKey, IncomePacket.CipherData[:ACD])
	if err != nil {
		return 0, err
	}
	packetData, err := aesCrypto.Decrypt(IncomePacket.CipherData[:ACD])
	if err != nil {
		return 0, err
	}
	defer memguard.WipeBytes(packetData)

	var PacketInfo Dto.GrpcIncomingPacketDetails
	err = f.Decoder.Decode(&PacketInfo, packetData)
	if err != nil {
		return 0, err
	}
	defer memguard.WipeBytes(PacketInfo.RsaKey)

	err = CheckTime(PacketInfo.TimeNow)
	if err != nil {
		return 0, err
	}

	err = f.CryptoValidate.CheckSign(DomainLevel.CheckSignKeyIncomingData{
		Sign: PacketInfo.Sign,
		Hash: func() []byte {
			hash := sha256.New()
			hash.Write(PacketInfo.RsaKey)
			return hash.Sum(nil)
		}(),
		MasterPublicKey: f.Keys.GetMasterPublicKey(),
	})
	if err != nil {
		return 0, err
	}

	f.TempKey.UpdateNewKey(PacketInfo.RsaKey)
	if err != nil {
		return 0, err
	}
	f.TempKey.UpdateOldKey()

	return PacketInfo.T1, nil
}
