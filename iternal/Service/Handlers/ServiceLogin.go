package Handlers

import (
	Dto2 "Kaban/iternal/Dto"
	"Kaban/iternal/Service/Helpers"
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func LoginCheck(password string, hashOfPassword string) error {
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

func LoginService(s Dto2.User, ctx context.Context) (string, string, error) {

	slog.Info("Func LoginService starts")
	app := *SetSettings()

	ctx, cancel := Helpers.ContextForDownloading(ctx)
	defer cancel()

	UnitId, password, err := app.Re.GetIdPassowrd(s.Email)
	if err != nil {
		slog.Error("Error in GetIdPassword")
		return "", "", err
	}

	err = LoginCheck(s.Password, password)
	if err != nil {
		slog.Error("func login 2", "err", err)
		return "", "", err
	}

	scrypt, err := app.Re.GeTScrypt(ctx, UnitId)
	if err != nil {
		slog.Error("Error in GeScrypt", "Err", err)
		return "", "", err
	}

	JWtToken, err := JWT(UnitId, scrypt)
	if err != nil {
		return "", "", err
	}
	RT, err := RFT(UnitId, scrypt)
	if err != nil {
		return "", "", err
	}

	slog.Info("Func LoginService ends")
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
