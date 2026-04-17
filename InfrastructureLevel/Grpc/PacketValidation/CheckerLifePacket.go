package PacketValidation

import (
	"time"
)

func (s *ValidatePacketData) CheckLifePacket(duration time.Time) bool {

	if time.Since(duration) > 5*time.Minute {
		return true
	}
	return false
}
