package DomainLevel

import "time"

type GrpcHandleData interface {
	SaveHash([]byte)
}
type PacketChecker interface {
	FindHash([32]byte) bool
	CheckLifePacket(time.Time) bool
	CheckSignature([]byte, []byte, []byte) error
}
type GrpcHandleRequest interface {
}
