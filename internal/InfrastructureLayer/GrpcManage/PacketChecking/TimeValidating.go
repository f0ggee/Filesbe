package PacketChecking

import (
	"errors"
	"time"
)

func (c PacketValidating) CheckTime(TimePacket time.Time) error {

	if time.Now().Before(TimePacket) {
		return errors.New("check time timeout")
	}
	return nil
}
