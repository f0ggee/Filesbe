package HandlingRequests

import (
	"Kaban/internal/Dto"
	"Kaban/internal/Service/Handlers"
	"crypto/rand"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/awnumar/memguard"
)

func (h HandlerGrpcRequest) CheckingGettingNewKey(Packet []byte) (time.Duration, error) {

	slog.Info("Start handling")
	PacketLook := Dto.GrpcOutComingPacketForSending{
		AesKeyData: nil,
		CipherData: nil,
	}

	err := json.Unmarshal(Packet, &PacketLook)
	if err != nil {
		slog.Error("Error while unmarshalling Packet", "Error", err.Error())
		return 0, err
	}

	slog.Debug("Packet Look Json", "Packet", PacketLook)

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

	err = h.ValidationPacket.CheckTime(PacketInfo.TimeNow)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(PacketData.Bytes(), &PacketInfo)
	if err != nil {
		slog.Info("PacketInfo", PacketInfo)
		slog.Error("Error while unmarshalling PacketInfo", "Error", err.Error())
		return 0, err
	}
	NewSavingRsa := memguard.NewBuffer(len(PacketInfo.RsaKey))
	NewSavingRsa.Copy(PacketInfo.RsaKey)
	memguard.WipeBytes(PacketInfo.RsaKey)

	err = h.CryptoValidate.CheckSignKey(PacketInfo.Sign, NewSavingRsa.Bytes(), Handlers.ControlPrivateKeyStruct.MasterServerPublicKeyBytes)
	if err != nil {
		return 0, err
	}

	Handlers.Keys.Mut.Lock()
	Handlers.Keys.NewPrivateKey = memguard.NewBuffer(NewSavingRsa.Size())
	Handlers.Keys.NewPrivateKey.Copy(NewSavingRsa.Bytes())
	Handlers.Keys.Mut.Unlock()

	Handlers.Keys.NewPrivateKey, err = memguard.NewBufferFromReader(rand.Reader, 32)

	if err != nil {
		slog.Error("Error while generating NewPrivateKey", "Error", err.Error())
	}

	slog.Info("Finish handling")

	return (PacketInfo.T1), nil
}
