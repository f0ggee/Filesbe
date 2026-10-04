package receive

import (
	"Kaban/internal/Protocol"
	"time"
)

func CheckTime(TimePacket time.Time) error {
	if time.Now().Before(TimePacket) {
		return Protocol.ErrorPacketTime
	}
	return nil
}
