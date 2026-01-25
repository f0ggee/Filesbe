package InfrastructureLayer

import (
	"Kaban/iternal/InfrastructureLayer/TokenInteraction"
	"Kaban/iternal/InfrastructureLayer/UserInteraction"
	"Kaban/iternal/Service/Connect_to_BD"
	"log/slog"
)

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

func SetSittingsTokenInteraction() *Generates {

	as := &TokenInteraction.ControlTokens{A: false}
	app := NewGenerateTokenEr(as)
	return app
}
