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

func RegisterService(de *Dto.HandlerRegister) (string, string, error) {

	app := *SetSettings()

	//Если пользователь существует, то тогда функция вернет кастомную ошибку - "person already exist"
	err := app.Re.CheckUser(de.Email)
	switch {
	case errors.Is(err, errors.New("person already exist")):
		return "", "", err

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

	//Интерфейс для создания пользователя, возвращает уникальный индификатор пользователя
	UnicIdUser, err := app.Re.CreateUser(de.Name, de.Email, HashPassword, ScryptKey)
	if err != nil {
		return "", "", err
	}

	TokenForJwt, err := JWT(UnicIdUser, ScryptKey)
	if err != nil {
		return "", "", err
	}
	TokenForRf, err := RFT(UnicIdUser, ScryptKey)
	if err != nil {
		return "", "", err
	}

	return TokenForJwt, TokenForRf, nil
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
		slog.Error("Error cant create password", err)
		return "", err
	}
	D4 := hex.EncodeToString(D3)
	return D4, nil
}
