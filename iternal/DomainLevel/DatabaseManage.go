package DomainLevel

import "context"

type ReadDb interface {
	GetIdPassword(string) (int, string, error)
	GetId(string, context.Context) (int, error)
	GetPassword(string, context.Context) (string, error)
}

type WriteDb interface {
	CreateUser(string, string, string, string) (int, error)
}
type CheckingDb interface {
	CheckerUser(string) error
}
