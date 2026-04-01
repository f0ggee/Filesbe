package HandlingRequests

import (
	"Kaban/internal/Dto"
	"Kaban/internal/Service/Handlers"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/awnumar/memguard"
	"golang.org/x/exp/rand"
)

func (h HandlerGrpcRequest) CheckingGettingNewKey(Packet []byte) (time.Duration, error) {

	slog.Info("Start handling", "ID", rand.Int())
	PacketLook := Dto.GrpcOutComingPacketForSending{
		AesKeyData: nil,
		CipherData: nil,
	}

	err := json.Unmarshal(Packet, &PacketLook)
	if err != nil {
		slog.Error("Error while unmarshalling Packet", "Error", err.Error())
		return 0, err
	}

	DecryptedAesKey, err := h.CryptoDecrypt.DecryptAesKey(Handlers.ControlPrivateKeyStruct.OurPrivateKeyIntoBytes, PacketLook.AesKeyData)
	if err != nil {
		return 0, err
	}

	PacketData := (h.CryptoDecrypt.DecryptPacket(DecryptedAesKey, PacketLook.CipherData))
	if PacketData == nil {
		return 0, errors.New("NewRsaKey error")
	}
	defer PacketData.Destroy()

	PacketInfo := Dto.GrpcIncomingPacketDetails{
		Sign:    nil,
		RsaKey:  nil,
		T1:      0,
		TimeNow: time.Now(),
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

	err = h.CryptoValidate.CheckSignKey(PacketInfo.Sign, Hash.Sum([]byte(nil)), Handlers.ControlPrivateKeyStruct.MasterServerPublicKeyBytes)
	if err != nil {
		return 0, err
	}
	h.Keys.UpdateKey(NewSavingRsa)

	slog.Info("Finish handling")

	return (PacketInfo.T1), nil
}
