package Handlers

import (
	"Kaban/iternal/Dto"
	"Kaban/iternal/Service/Helpers"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"

	"golang.org/x/crypto/scrypt"
)

func (sa *HandlerPackCollect) RegisterService(de *Dto.HandlerRegister) (string, string, error) {

	//app := *InfrastructureLayer.SetSettings()
	//manageTokenInteraction := *InfrastructureLayer.SetSittingsTokenInteraction()

	//Если пользователь существует, то тогда функция вернет кастомную ошибку - "person already exist"
	err := sa.S.Database.CheckUser(de.Email)
	switch {
	case errors.Is(err, errors.New("person already exist")):
		return "", "", errors.New("person already exist")

	case err != nil:
		return "", "", err
	}
	HashPassword, err := Helpers.HashPassowrd(de.Password)
	if err != nil {
		slog.Error("Err generate a password-scrypt", "err", err)
		return "", "", err
	}

	ScryptKey, err := GenerateScrypt(de, err)
	if err != nil {
		return "", "", err
	}

	UnicIdUser, err := sa.S.Database.CreateUser(de.Name, de.Email, HashPassword, ScryptKey)
	if err != nil {
		return "", "", err
	}

	DataCollected := sa.S.TokenImpl.CollectDataForTokens(UnicIdUser)

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

	return JwtToken, RefreshToken, nil
}

func GenerateScrypt(de *Dto.HandlerRegister, err error) (string, error) {
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		slog.Error("Err in fill slice", "Err", err)
		return "", err
	}
	D3, err := scrypt.Key([]byte(de.Password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		slog.Error("Error cant create password", err.Error())
		return "", err
	}
	D4 := hex.EncodeToString(D3)
	return D4, nil
}
