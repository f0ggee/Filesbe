package InfrastructureLayer

import (
	"Kaban/iternal/DomainLevel"
	"Kaban/iternal/InfrastructureLayer/KeyInteration"
)

type EncryptionKeyInteration struct {
	Choose DomainLevel.KeyInteraction
}

func NewSetEncryptionKeyInteration(Connect DomainLevel.KeyInteraction) *EncryptionKeyInteration {
	return &EncryptionKeyInteration{Choose: Connect}
}

func ConnectToEncryptKey() *EncryptionKeyInteration {

	apps := &KeyInteration.EncryptionKey{}
	saz := NewSetEncryptionKeyInteration(apps)
	return saz
}
