package InfrastructureLayer

import (
	"context"
)

type RepositorysForUsingDb interface {
	GeTScrypt(context.Context, int) (string, error)
	GetIdPassowrd(string) (int, string, error)
	CheckUser(string) error
	CreateUser(string, string, string, string) (int, error)
}

type GenerateToken interface {
	GenerateJWT(int, string)
}
type RepositorysTest interface {
	GetIdPassowrdTest(string) (int, string, error)

	GeTScryptTest(context.Context, int) (string, error)
}
