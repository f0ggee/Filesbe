package HandlingRequests

import (
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/Crypto/Decription"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/awnumar/memguard"
	"golang.org/x/exp/rand"
)

func (h HandlerGrpcRequest) CheckingGettingNewKey(Packet []byte) (time.Duration, error) {

	Id := rand.Int()
	slog.Info("Start accepting the new key", slog.Int("ID", Id))
	PacketLook := Dto.GrpcOutComingPacketForSending{
		AesKeyData: nil,
		CipherData: nil,
	}

	err := json.Unmarshal(Packet, &PacketLook)
	if err != nil {
		slog.Error("Error while unmarshalling Packet", "Error", err.Error())
		return 0, err
	}

	DecryptedAesKey, err := h.CryptoDecrypt.DecryptAesKey(h.Keys.GetOurKey(), PacketLook.AesKeyData)
	if err != nil {
		return 0, err
	}

	PacketData := h.CryptoDecrypt.DecryptPacket(DecryptedAesKey, PacketLook.CipherData)
	if PacketData == nil {
		return 0, errors.New("NewRsaKey error")
	}
	defer PacketData.Destroy()
	PacketInfo := &Dto.GrpcIncomingPacketDetails{
		Sign:   nil,
		RsaKey: nil,
		T1:     0,
	}
	err = json.Unmarshal(PacketData.Bytes(), &PacketInfo)
	if err != nil {
		slog.Error("Error while unmarshalling PacketInfo", "Error", err.Error())
		return 0, err
	}
	err = h.ValidationPacket.CheckTime(PacketInfo.TimeNow)
	if err != nil {
		return 0, err
	}

	NewSavingRsa := memguard.NewBuffer(len(PacketInfo.RsaKey))
	NewSavingRsa.Copy(PacketInfo.RsaKey)
	memguard.WipeBytes(PacketInfo.RsaKey)
	defer NewSavingRsa.Destroy()

	Hash := sha256.New()

	Hash.Write(NewSavingRsa.Bytes())

	err = h.CryptoValidate.CheckSignKey(PacketInfo.Sign, Hash.Sum([]byte(nil)), h.Keys.GetMasterKey())
	if err != nil {
		return 0, err
	}

	h.Keys.UpdateKey(NewSavingRsa)
	h.Keys.UpdateOldKey()

	slog.Info("Finish accepting the new key", slog.Int("Id", Id))

	return PacketInfo.T1, nil
}

func Ee2(sa *Decription.DecryptionData, key, data []byte) error {

	_, err := sa.DecryptAesKey(key, data)
	if err != nil {
		return err
	}
	return nil
}
