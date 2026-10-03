package Protocol

import "errors"

var (
	ErrorGetKey      = errors.New("can't get a new key")
	ErrorPrepareData = errors.New("can't prepare data")
	ErrorPacketTime  = errors.New("packet's time isn't correct")
)
