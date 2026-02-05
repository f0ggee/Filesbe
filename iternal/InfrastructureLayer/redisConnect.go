package InfrastructureLayer

import (
	"Kaban/iternal/DomainLevel"
	"Kaban/iternal/InfrastructureLayer/RedisInteration"
)

type RedisInterationLayer struct {
	Ras DomainLevel.RedisInteration
}

func NewRedisInterationLayer(rep DomainLevel.RedisInteration) *RedisInterationLayer {
	return &RedisInterationLayer{Ras: rep}
}

func NewSetRedisConnect() *RedisInterationLayer {

	app := &RedisInteration.RedisInterationLayer{}
	sa := NewRedisInterationLayer(app)
	return sa
}
