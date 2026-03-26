package DomainLevel

import "github.com/awnumar/memguard"

type DataConvert interface {
	JsonConverter(any) ([]byte, error)
}

type KeysManager interface {
	UpdateKey(key *memguard.LockedBuffer)
	GetKey() []byte
	GetOldKey() []byte
	UpdateOldKey()
}
