package GrpcManage

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/KeysManager"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"log/slog"
	"time"

	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"github.com/awnumar/memguard"
)

type HandlerGrpcRequest struct {
	NewValidating
	NewKeys
	NewDecrypt
	NewEncrypt
	NewParser
}

func GetNewHandlerGrpcRequest(newValidating NewValidating, newKeys NewKeys, newDecrypt NewDecrypt, newEncrypt NewEncrypt, newParser NewParser) *HandlerGrpcRequest {
	return &HandlerGrpcRequest{NewValidating: newValidating, NewKeys: newKeys, NewDecrypt: newDecrypt, NewEncrypt: newEncrypt, NewParser: newParser}
}

type NewKeys struct {
	Keys       KeysManager.Keys
	ServerKeys DomainLevel.NewSetKeys
}
type NewValidating struct {
	CryptoValidate DomainLevel.CryptoValidating
}
type NewDecrypt struct {
	CryptoDecrypt DomainLevel.Decryption
}
type NewEncrypt struct {
	CryptoEncrypt DomainLevel.Encryption
}
type NewParser struct {
	Parse RepoParsers.Decode
}

func (h HandlerGrpcRequest) CheckingGettingNewKey(Packet []byte) (time.Duration, error) {

	Id := rand.Text()
	slog.Info("Func CheckingGettingNewKey: Start accepting the new key", slog.String("ID", Id))
	PacketLook := Dto.GrpcOutComingPacketForSending{
		AesKeyData: nil,
		CipherData: nil,
	}

	err := h.Parse.JsonDecodeMarshall(&PacketLook, Packet)
	if err != nil {
		slog.Error("Error while unmarshalling Packet", "Error", err.Error())
		return 0, err
	}

	DecryptedAesKey, err := h.CryptoDecrypt.DecryptAesKey(h.ServerKeys.GerOurPrivateKey(), PacketLook.AesKeyData)
	if err != nil {
		return 0, err
	}

	PacketData := h.CryptoDecrypt.DecryptPacket(DecryptedAesKey, PacketLook.CipherData)
	if PacketData == nil {
		return 0, errors.New("NewRsaKey error")
	}
	defer PacketData.Destroy()

	PacketInfo := &Dto.GrpcIncomingPacketDetails{
		Sign:    nil,
		RsaKey:  nil,
		T1:      0,
		TimeNow: time.Now(),
	}
	err = h.Parse.JsonDecodeMarshall(&PacketInfo, PacketData.Bytes())
	if err != nil {
		slog.Error("CheckingGettingNewKey;Error while unmarshalling PacketInfo", "Error", err.Error())
		return 0, err
	}
	err = h.setCheckTime(PacketInfo.TimeNow)
	if err != nil {
		return 0, err
	}

	NewSavingRsa, err := h.getSaveKey(&PacketInfo.RsaKey)
	if err != nil {
		return DomainLevel.DefaultErrorTime, nil
	}
	defer NewSavingRsa.Destroy()

	Hash := sha256.New()
	Hash.Write(NewSavingRsa.Bytes())
	err = h.CryptoValidate.CheckSignKey(DomainLevel.CheckSignKeyIncomingData{
		Sign:            PacketInfo.Sign,
		Hash:            Hash.Sum(nil),
		MasterPublicKey: h.ServerKeys.GetMasterPublicKey(),
	})
	if err != nil {
		return 0, err
	}

	h.Keys.UpdateOldKey()
	h.Keys.UpdateKey(NewSavingRsa)

	slog.Info("Finish accepting the new key", slog.Group("Data",
		slog.String("Id", Id),
		slog.Duration("Time for next swapping", PacketInfo.T1)))

	return PacketInfo.T1, nil
}

func (h *HandlerGrpcRequest) setCheckTime(TimePacket time.Time) error {
	if time.Now().Before(TimePacket) {
		return errors.New("check time timeout")
	}
	return nil
}
func (h HandlerGrpcRequest) getSaveKey(data *[]byte) (*memguard.LockedBuffer, error) {
	if data == nil {
		return nil, errors.New(DomainLevel.ErrorDataNil)
	}
	NewSavingRsa := memguard.NewBuffer(len(*data))
	NewSavingRsa.Copy(*data)
	memguard.WipeBytes(*data)
	return NewSavingRsa, nil
}
