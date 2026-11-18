package Application

import (
	"Kaban/internal/DomainLevel"
	"context"
	"io"
)

type DownloadApplication interface {
	DownloadApplication(ctx context.Context, fileName string, dst io.Writer) error
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
