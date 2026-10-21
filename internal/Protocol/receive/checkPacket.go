package receive

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"Kaban/internal/Protocol"
	"time"

	"github.com/awnumar/memguard"
)

type GetNewKey struct {
	Decoder  DomainLevel.Decoder
	Aes      DomainLevel.CryptoMaker
	Rsa      DomainLevel.CryptoMaker
	Validate DomainLevel.SignatureMaker
}

// GetPacketData decrypts and cheks a packet and then returns data.
func (f GetNewKey) GetPacketData(packet []byte) (time.Duration, error) {
	var IncomePacket Protocol.OutComePacket
	err := f.Decoder.Decode(&IncomePacket, packet)
	if err != nil {
		return 0, err
	}

	var ACD = f.Aes.GetRequiredOverheadSize()

	rsaDecrypter, err := f.Rsa.MakeCrypto(DomainLevel.ServerKeys.GerOurPrivateKey(), nil)
	if err != nil {
		return 0, err
	}

	decryptedAesKey, err := rsaDecrypter.Decrypt(IncomePacket.AesKeyData)
	if err != nil {
		return 0, err
	}

	aesCrypto, err := f.Aes.MakeCrypto(decryptedAesKey, IncomePacket.CipherData[:ACD])
	if err != nil {
		return 0, err
	}
	packetData, err := aesCrypto.Decrypt(IncomePacket.CipherData[:ACD])
	if err != nil {
		return 0, err
	}
	defer memguard.WipeBytes(packetData)

	var PacketInfo Protocol.IncomePacketDetails
	err = f.Decoder.Decode(&PacketInfo, packetData)
	if err != nil {
		return 0, err
	}
	defer memguard.WipeBytes(PacketInfo.RsaKey)

	err = CheckTime(PacketInfo.TimeNow)
	if err != nil {
		return 0, err
	}

	signMaker, err := f.Validate.MakeSigner(DomainLevel.ServerKeys.GetMasterPublicKey())
	if err != nil {
		return 0, err
	}
	err = signMaker.CheckSign(PacketInfo.Sign)
	if err != nil {
		return 0, err
	}
	RepoEncrypterKeys.KeysControl.UpdateNewKey(PacketInfo.RsaKey)
	RepoEncrypterKeys.KeysControl.UpdateOldKey()
	return PacketInfo.T1, nil
}
