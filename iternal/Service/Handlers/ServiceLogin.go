package Handlers

import (
	Dto2 "Kaban/iternal/Dto"
	"Kaban/iternal/Service/Helpers"
	"context"
	"crypto/rand"
	"log/slog"
	"time"

	"Kaban/iternal/Dto"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func PasswordCheck(password string, hashOfPassword string) error {
	slog.Info("Password checking starts")
	err := bcrypt.CompareHashAndPassword([]byte(hashOfPassword), []byte(password))
	if err != nil {
		slog.Error("Error while checking the password", "Error", err.Error())
		return err

	}
	slog.Info("Password ends")
	return nil

}

func CollectData(id int) Dto.JwtCustomStruct {
	ds := Dto.JwtCustomStruct{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Kabaner",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			ID:        rand.Text(),
		},
	}

	return ds
}

func (sa *HandlerPackCollect) LoginService(s Dto2.User, ctx context.Context) (string, string, error) {

	slog.Info("Func LoginService starts")
	//app := *InfrastructureLayer.SetSettings()
	//
	//ManageTokenApp := *InfrastructureLayer.SetSittingsTokenInteraction()

	ctx, cancel := Helpers.ContextForDownloading(ctx)
	defer cancel()

	Id, password, err := sa.S.Database.GetIdPassword(s.Email)
	//Id, password, err := app.Re.GetIdPassword(s.Email)
	if err != nil {
		slog.Error("Error in GetIdPassword", "error", err)
		return "", "", err
	}

	err = PasswordCheck(s.Password, password)
	if err != nil {
		slog.Error("func login 2", "err", err)
		return "", "", err
	}

	DataCollected := CollectData(Id)

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
