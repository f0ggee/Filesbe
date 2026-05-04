package Decription

import (
	"crypto/aes"
	"crypto/cipher"
	"log/slog"

	"github.com/awnumar/memguard"
)

func (d DecryptionData) DecryptPacket(aesKey []byte, plainText []byte) *memguard.LockedBuffer {
	slog.Info("Func DecryptPacket: Start decrypting a packet")
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		slog.Error("Func DecryptPacket:Error create new aes block", "Error", err.Error())
		return nil
	}
	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		slog.Error("Func DecryptPacket: Error create new gcm", "Error", err.Error())
		return nil
	}
	sa, err := gcm.Open(nil, plainText[:gcm.NonceSize()], plainText[gcm.NonceSize():], nil)

	if err != nil {
		slog.Error("Func DecryptPacket: Error decrypt packet", "Error", err.Error())
		return nil
	}
	defer memguard.WipeBytes(sa)
	saz := memguard.NewBufferFromBytes(sa)
	slog.Info("Func DecryptPacket: Finish decrypting a packet")
	return saz
}
