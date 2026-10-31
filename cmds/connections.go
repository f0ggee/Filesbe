package cmds

import (
	"Kaban/internal/InfrastructureLayer/DatabaseControl"
	"Kaban/internal/InfrastructureLayer/Grpc"
	"Kaban/internal/InfrastructureLayer/RedisInteration"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type Conns struct {
	Conn        *redis.Client
	Db          *pgxpool.Pool
	GrpcConnect *grpc.ClientConn
}

func (c Conns) CloseConnections() {
	c.GrpcConnect.Close()
	c.Conn.Close()
	c.Db.Close()
}

func GetConnections() Conns {
	conn := RedisInteration.GetRedisConn()

	db, err := DatabaseControl.Connect()
	if err != nil {
		panic(err)
	}
	grpcConnect, err := Grpc.MakeGrpcConn()
	if err != nil {
		panic(err)
	}

	return Conns{
		Conn:        conn,
		Db:          db,
		GrpcConnect: grpcConnect,
	}
}
