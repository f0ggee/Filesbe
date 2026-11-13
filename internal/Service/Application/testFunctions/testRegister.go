package testFunctions

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Drivers/Crypto"
	"Kaban/internal/Drivers/DatabaseControl"
	"Kaban/internal/Drivers/Tokens"
	"Kaban/internal/Service/Application"
	"errors"
)

func GetTestRegister() Application.RegisterApplication {
	return &Application.NewRegisterApplication{
		CheckingDb: DatabaseControl.NewCheckerTest(),
		WriterDb:   DatabaseControl.NewWriterTest(),
		Generator:  Crypto.GetNewGenerating(),
		Rft:        Tokens.NewJwtToken(),
		Jwt:        Tokens.NewRfToken(),
	}
}

type RegisterError struct {
}

func GetRegisterError() *RegisterError {
	return &RegisterError{}
}

func (r RegisterError) RegisterApplication(register Application.InputRegister) DomainLevel.RegisterApplicationOutComingData {

	return DomainLevel.RegisterApplicationOutComingData{
		Err: errors.New("error"),
	}
}
