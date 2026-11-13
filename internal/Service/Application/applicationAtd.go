package Application

import (
	"Kaban/internal/DomainLevel"
	"context"
)

type DownloadApplication interface {
	DownloadApplication(context.Context, string) error
}

type DownloadEncryptApplication interface {
	DownloadEncryptApplication(DownloadEncryptIncomeData) error
}
type InputLogin struct {
	Ctx      context.Context
	Email    string
	Password []byte
}
type LoginApplication interface {
	LoginApplication(InputLogin) DomainLevel.LoginApplicationOutComingData
}
type InputRegister struct {
	Ctx      context.Context
	Name     string
	Email    string
	Password []byte
}
type RegisterApplication interface {
	RegisterApplication(InputRegister) DomainLevel.RegisterApplicationOutComingData
}
type UploadApplication interface {
	UploadApplication(UploadIncomeData) (string, error)
}
type UploaderEncryptApplication interface {
	UploadEncryptApplication(UploadEncryptIncomeData) (string, error)
}
