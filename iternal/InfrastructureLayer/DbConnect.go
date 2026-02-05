package InfrastructureLayer

import (
	"Kaban/iternal/DomainLevel"
	"Kaban/iternal/InfrastructureLayer/UserInteraction"
	"Kaban/iternal/Service/Connect_to_BD"
	"log/slog"
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
	db, err := Connect_to_BD.Connect()
	if err != nil {
		slog.Error("Err_from_register 1 ", err)
		return nil
	}
	S := &UserInteraction.DB{Db: db}
	app := NewUserService(S)
	return app
}
func NewUserServiceTest(Reps DomainLevel.RepositorysTest) *ConnectToBdTests {
	return &ConnectToBdTests{Res: Reps}
}
