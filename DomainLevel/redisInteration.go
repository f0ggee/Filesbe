package DomainLevel

type RedisUse interface {
	SendData([]byte, string) error
}
