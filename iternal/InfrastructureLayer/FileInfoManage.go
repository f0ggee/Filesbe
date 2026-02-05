package InfrastructureLayer

import (
	"Kaban/iternal/DomainLevel"
	"Kaban/iternal/InfrastructureLayer/FileKeyInteration"
)

func ConnectKeyControl() *ControlAccessKeys {
	apps := &FileKeyInteration.FileInfoController{}
	sa := NewSetKeyController(apps)
	return sa
}

type ControlAccessKeys struct {
	Key DomainLevel.FileInfo
}

func NewSetKeyController(r DomainLevel.FileInfo) *ControlAccessKeys {
	return &ControlAccessKeys{Key: r}
}
