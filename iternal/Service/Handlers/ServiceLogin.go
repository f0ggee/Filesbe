package Handlers

import (
	Dto2 "Kaban/iternal/Dto"
	"Kaban/iternal/Service/Helpers"
	"context"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

func PasswordCheck(password string, hashOfPassword string) error {
	slog.Info("Func LoginCheck starts")
	err := bcrypt.CompareHashAndPassword([]byte(hashOfPassword), []byte(password))
	if err != nil {
		slog.Error("Err", "err", err)
		return err

	}
	slog.Info("Func LoginCheck ends")
	return nil

}

//func SetSettingsTest(db string) *InfrastructureLayer2.ConnectToBdTests {
//	S := &InfrastructureLayer2.DbForTests{DbTest: db}
//	app := InfrastructureLayer2.NewUserServiceTest(S)
//	return app
//}

func (sa *HandlerPackCollect) LoginService(s Dto2.User, ctx context.Context) (string, string, error) {

	slog.Info("Func LoginService starts")
	//app := *InfrastructureLayer.SetSettings()
	//
	//ManageTokenApp := *InfrastructureLayer.SetSittingsTokenInteraction()

	ctx, cancel := Helpers.ContextForDownloading(ctx)
	defer cancel()

	Id, password, err := sa.S.Database.GetIdPassowrd(s.Email)
	//Id, password, err := app.Re.GetIdPassowrd(s.Email)
	if err != nil {
		slog.Error("Error in GetIdPassword", "error", err)
		return "", "", err
	}

	err = PasswordCheck(s.Password, password)
	if err != nil {
		slog.Error("func login 2", "err", err)
		return "", "", err
	}

	DataCollected := sa.S.TokenImpl.CollectDataForTokens(Id)

	RefreshToken, err := sa.S.Tokens.GenerateRT(DataCollected)
	if err != nil {
		slog.Error("func login 3", "err", err)
		return "", "", err
	}
	JwtToken, err := sa.S.Tokens.GenerateJWT(DataCollected)
	if err != nil {
		slog.Error("func login 4", "err", err)
		return "", "", err
	}

	slog.Info("Func LoginService ends")
	return JwtToken, RefreshToken, nil

}
