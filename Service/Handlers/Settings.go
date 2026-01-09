package Handlers

import (
	"Kaban/InfrastructureLayer"
	"Kaban/Service/Connect_to_BD"
	"log/slog"
)

func SetSettings() *InfrastructureLayer.ConnectToBd {
	db, err := Connect_to_BD.Connect()
	if err != nil {
		slog.Error("Err_from_register 1 ", err)
		return nil
	}

	S := &InfrastructureLayer.DB{Db: db}
	app := InfrastructureLayer.NewUserService(S)
	return app
}
