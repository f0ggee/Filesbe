package InfrastructureLayer

import "Kaban/iternal/InfrastructureLayer/FileKeyInteration"

func ConnectKeyControl() *ControlAccessKeys {
	apps := &FileKeyInteration.FileInfoController{}
	sa := NewSetKeyController(apps)
	return sa
}
