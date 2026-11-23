package testFunctions

import (
	"Kaban/internal/Drivers/Crypto"
	"Kaban/internal/Drivers/DatabaseControl"
	"Kaban/internal/Drivers/Tokens"
	"Kaban/internal/Service/Application"
)

func NewTestLogin() Application.LoginApplication {
	return &Application.NewLogin{
		ReaderDatabase: DatabaseControl.NewReaderTest(),
		Validate:       Crypto.NewCheckingTest(),
		Rf:             Tokens.NewRfToken(),
		Jwt:            Tokens.NewJwtToken(),
	}
}
