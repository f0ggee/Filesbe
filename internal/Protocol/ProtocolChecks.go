package Protocol

import (
	"time"
)

func CheckTime(TimePacket time.Time) error {
	if time.Now().Before(TimePacket) {
		return ErrorPacketTime
	}
	return nil
}
