package InfrastructureLayer

import "Kaban/iternal/DomainLevel"

type ControlAccessKeys struct {
	Key DomainLevel.FileInfo
}

func NewSetKeyController(r DomainLevel.FileInfo) *ControlAccessKeys {
	return &ControlAccessKeys{Key: r}
}
