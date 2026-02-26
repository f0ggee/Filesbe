package InfrastructureLayer

import (
	"Kaban/iternal/DomainLevel"
	"Kaban/iternal/InfrastructureLayer/UserInteraction"
)

type ConnectToBd struct {
	Re DomainLevel.UserServer
}

func NewUserService(Rep DomainLevel.UserServer) *ConnectToBd {
	return &ConnectToBd{Re: Rep}
}

type ConnectToBdTests struct {
	Res DomainLevel.RepositorysTest
}

func SetSettings() *ConnectToBd {

	S := &UserInteraction.DB{Db: nil}
	app := NewUserService(S)
	return app
}
func NewUserServiceTest(Reps DomainLevel.RepositorysTest) *ConnectToBdTests {
	return &ConnectToBdTests{Res: Reps}
}
