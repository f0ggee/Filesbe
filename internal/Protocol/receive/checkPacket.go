package receive

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"crypto/sha256"
	"time"

	"github.com/awnumar/memguard"
)

type GetNewKey struct {
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

	rsaDecrypter, err := f.rsa.MakeCrypto(DomainLevel.ServerKeys.GerOurPrivateKey(), nil)
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
		MasterPublicKey: DomainLevel.ServerKeys.GetMasterPublicKey(),
	})
	if err != nil {
		return 0, err
	}

	f.TempKey.UpdateNewKey(PacketInfo.RsaKey)
	f.TempKey.UpdateOldKey()

	return PacketInfo.T1, nil
}
