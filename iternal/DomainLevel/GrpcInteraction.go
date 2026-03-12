package DomainLevel

type GrpcInteraction interface {
	SendRequestGrpc([]byte) ([]byte, error)
	SayHi() string
}

type GrpcDataManage interface {
	GenerateSignature(message []byte, key []byte) ([]byte, error)
	SayHI() string
}
