package InfrastructureLayer

type ConnectToBd struct {
	Re RepositorysForUsingDb
}

func NewUserService(Rep RepositorysForUsingDb) *ConnectToBd {
	return &ConnectToBd{Re: Rep}
}

type ConnectToBdTests struct {
	Res RepositorysTest
}

func NewUserServiceTest(Reps RepositorysTest) *ConnectToBdTests {
	return &ConnectToBdTests{Res: Reps}
}
