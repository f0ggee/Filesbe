package InfrastructureLayer

import "Kaban/iternal/DomainLevel"

type ConnectToBd struct {
	Re DomainLevel.UserServer
}

func NewUserService(Rep DomainLevel.UserServer) *ConnectToBd {
	return &ConnectToBd{Re: Rep}
}

type ConnectToBdTests struct {
	Res DomainLevel.RepositorysTest
}

func NewUserServiceTest(Reps DomainLevel.RepositorysTest) *ConnectToBdTests {
	return &ConnectToBdTests{Res: Reps}
}
