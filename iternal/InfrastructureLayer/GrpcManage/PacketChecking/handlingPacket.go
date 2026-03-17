package PacketChecking

import (
	"Kaban/iternal/Dto"
	"Kaban/iternal/Service/Handlers"
	"crypto/rand"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/awnumar/memguard"
)

func (v *Validating) Handle(Packet []byte) (time.Duration, error) {

	slog.Info("Start handling")
	PacketLook := Dto.GrpcSendingData{
		AesKeyData: nil,
		CipherData: nil,
	}

	err := json.Unmarshal(Packet, &PacketLook)
	if err != nil {
		slog.Error("Error while unmarshalling Packet", "Error", err.Error())
		return 0, err
	}

	slog.Debug("Packet Look Json", "Packet", PacketLook)

	DecryptedAesKey, err := v.Decrypter.DecryptAesKey(Handlers.ControlPrivateKeyStruct.OurPrivateKeyIntoBytes, PacketLook.AesKeyData)
	if err != nil {
		return 0, err
	}

	PacketData := (v.Decrypter.DecryptPacket(DecryptedAesKey, PacketLook.CipherData))
	if PacketData == nil {
		return 0, errors.New("NewRsaKey error")
	}
	defer PacketData.Destroy()

	PacketInfo := Dto.GrpcOutcomingDataPre{
		Sign:   nil,
		RsaKey: nil,
		T1:     0,
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

	err = v.CheckSign(PacketInfo.Sign, NewSavingRsa.Bytes(), Handlers.ControlPrivateKeyStruct.MasterServerPublicKeyBytes)
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
