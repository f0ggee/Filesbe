package KeyInteration

import (
	"crypto/aes"
	"crypto/cipher"
	"log/slog"

	"github.com/awnumar/memguard"
)

func (*EncryptionKey) DecryptPacket(aesKey []byte, plainText []byte) *memguard.LockedBuffer {
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		slog.Error("Error create new aes block", "Error", err.Error())
		return nil
	}
	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		slog.Error("Error create new gcm", "Error", err.Error())
		return nil
	}
	sa, err := gcm.Open(nil, plainText[:gcm.NonceSize()], plainText[gcm.NonceSize():], nil)

	if err != nil {
		slog.Error("Error decrypt packet", "Error", err.Error())
		return nil
	}
	defer memguard.WipeBytes(sa)
	saz := memguard.NewBufferFromBytes(sa)
	return saz
}
