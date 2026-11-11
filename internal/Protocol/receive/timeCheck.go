package receive

import (
	"Kaban/internal/Protocol"
	"time"
)

var packetLive = 2 * time.Minute

func CheckTime(TimePacket time.Time) error {

	if time.Now().Sub(TimePacket) > packetLive {
		return Protocol.ErrorPacketTime
	}
	return nil
}
