package DomainLevel

import "context"

type OutComingLoginData struct {
	Id           int32
	HashPassword string
	Err          error
}
type CreateUserIncomingData struct {
	Name         string
	Email        string
	HashPassword string
	Ctx          context.Context
}

type ReadDb interface {
	LoginData(context.Context, string) OutComingLoginData
}

type WriteDb interface {
	CreateUser(CreateUserIncomingData) (int32, error)
}
type CheckingDb interface {
	CheckerUser(context.Context, string) error
}
