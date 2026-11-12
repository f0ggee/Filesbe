package Protocol

import "time"

const (
	DefaultErrorTime = 12 * time.Hour
	ErrorTimePacket  = "the packet time isn't correct"
)

type InEnvelope struct {
	Sign    []byte        `json:"Sign"`
	RsaKey  []byte        `json:"MasterPublicKey"`
	T1      time.Duration `json:"T1"`
	TimeNow time.Time     `json:"TimeNow"`
}
type OutPacket struct {
	AesKeyData []byte `json:"aes_key_data"`
	CipherData []byte `json:"cipher_data"`
}

type MasterServerEnvelope struct {
	Time             time.Time
	ServerName       []byte
	SignedServerName []byte
}
