package DomainLevel

import "context"

type ReadDb interface {
	LoginData(string, context.Context) (int, string, error)
}

type WriteDb interface {
	CreateUser(string, string, string, context.Context) (int, error)
}
type CheckingDb interface {
	CheckerUser(string, context.Context) error
}
