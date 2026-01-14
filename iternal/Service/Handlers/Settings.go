package Handlers

import (
	InfrastructureLayer2 "Kaban/iternal/InfrastructureLayer"
	"Kaban/iternal/Service/Connect_to_BD"
	"log/slog"
)

func SetSettings() *InfrastructureLayer2.ConnectToBd {
	db, err := Connect_to_BD.Connect()
	if err != nil {
		slog.Error("Err_from_register 1 ", err)
		return nil
	}

	S := &InfrastructureLayer2.DB{Db: db}
	app := InfrastructureLayer2.NewUserService(S)
	return app
}
