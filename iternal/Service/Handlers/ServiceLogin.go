package Handlers

import (
	Dto2 "Kaban/iternal/Dto"
	"Kaban/iternal/Service/Uttiltesss"
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func LoginCheck(password string, hash_of_password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash_of_password), []byte(password))
	if err != nil {
		slog.Error("Err", "err", err)
		return err

	}
	return nil

}

//func SetSettingsTest(db string) *InfrastructureLayer2.ConnectToBdTests {
//	S := &InfrastructureLayer2.DbForTests{DbTest: db}
//	app := InfrastructureLayer2.NewUserServiceTest(S)
//	return app
//}

func LoginService(s Dto2.User, ctx context.Context) (string, string, error) {

	app := *SetSettings()

	ctx, cancel := Uttiltesss.Contexte(ctx)
	defer cancel()

	UnicId, password, err := app.Re.GetIdPassowrd(s.Email)
	if err != nil {
		slog.Error("Error in GetIdPassword")
		return "", "", err
	}

	err = LoginCheck(s.Password, password)
	if err != nil {
		slog.Error("func login 2", "err", err)
		return "", "", err
	}

	scrypt, err := app.Re.GeTScrypt(ctx, UnicId)
	if err != nil {
		slog.Error("Error in GeScrypt", "Err", err)
		return "", "", err
	}

	JWtToken, err := JWT(UnicId, scrypt)
	if err != nil {
		return "", "", err
	}
	RT, err := RFT(UnicId, scrypt)
	if err != nil {
		return "", "", err
	}

	return JWtToken, RT, nil

}

func JWT(UnicId int, scrypt string) (string, error) {

	claims := Dto2.MyCustomCookie{
		UserID: UnicId,
		Sc:     scrypt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Admin",
		},
	}
	token1 := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	sd := []byte(os.Getenv("KEYFORJWT"))

	tokenString, err := token1.SignedString(sd)
	if err != nil {
		slog.Error("Error sign cookie jwt", err)
		return "", err
	}

	return tokenString, nil
}

func RFT(UnicId int, scrypt string) (string, error) {
	claims := Dto2.MyCustomCookie{
		UserID: UnicId,
		Sc:     scrypt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Admin",
		},
	}
	token1 := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	sd := []byte(os.Getenv("KEYFORJWT"))
	tokenString, err := token1.SignedString(sd)
	if err != nil {
		slog.Error("Error sign cookie", err)
		return "", err
	}

	return tokenString, nil
}
