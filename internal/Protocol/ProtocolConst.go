package Protocol

import "time"

const (
	DefaultErrorTime = 12 * time.Hour
	ErrorTimePacket  = "the packet time isn't correct"
)

type IncomePacketDetails struct {
	Sign    []byte        `json:"Sign"`
	RsaKey  []byte        `json:"MasterPublicKey"`
	T1      time.Duration `json:"T1"`
	TimeNow time.Time     `json:"TimeNow"`
}
type OutComePacket struct {
	AesKeyData []byte `json:"aes_key_data"`
	CipherData []byte `json:"cipher_data"`
}

type OutPacketDetails struct {
	Time             time.Time
	ServerName       []byte
	SignedServerName []byte
}
