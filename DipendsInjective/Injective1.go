package DipendsInjective

import "MasterServer_/DomainLevel"

type RsaKeyManipulationWithRsaAndMemory struct {
	KeyAndMemory DomainLevel.KeyManage
	Key          DomainLevel.RsaKeyManipulation
}

func NewRsaKeyManipulationWithRsaAndMemory(keyAndMemory DomainLevel.KeyManage, rsaKey DomainLevel.RsaKeyManipulation) *RsaKeyManipulationWithRsaAndMemory {
	return &RsaKeyManipulationWithRsaAndMemory{KeyAndMemory: keyAndMemory, Key: rsaKey}
}
