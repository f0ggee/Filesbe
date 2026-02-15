package InfrastructureLayer

import (
	"Kaban/iternal/DomainLevel"
	"Kaban/iternal/InfrastructureLayer/encryptionKeyInteration"
)

type EncryptionKeyInteration struct {
	Choose DomainLevel.EncryptionKeyInteration
}

func NewSetEncryptionKeyInteration(Connect DomainLevel.EncryptionKeyInteration) *EncryptionKeyInteration {
	return &EncryptionKeyInteration{Choose: Connect}
}

func ConnectToEncryptKey() *EncryptionKeyInteration {

	apps := &encryptionKeyInteration.EncryptionKey{}
	saz := NewSetEncryptionKeyInteration(apps)
	return saz
}
