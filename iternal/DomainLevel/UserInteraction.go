package DomainLevel

import (
	"context"
)

type UserServer interface {
	GeTScrypt(context.Context, int) (string, error)
	GetIdPassowrd(string) (int, string, error)
	CheckUser(string) error
	CreateUser(string, string, string, string) (int, error)
}
