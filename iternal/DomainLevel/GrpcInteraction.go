package DomainLevel

type GrpcInteraction interface {
	SendData([]byte) ([]byte, error)
}
