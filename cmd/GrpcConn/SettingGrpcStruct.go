package GrpcConn

type GrpcConnection struct {
	Network string
	Address string
}

func GetGrpcConn(network string, address string) *GrpcConnection {
	return &GrpcConnection{
		Network: network,
		Address: address,
	}
}
