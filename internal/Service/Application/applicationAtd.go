package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"context"
)

type DownloadApplication interface {
	DownloadApplication(context.Context, string) error
}

type DownloadEncryptApplication interface {
	DownloadEncryptApplication(DownloadEncryptIncomeData) error
}

type LoginApplication interface {
	LoginApplication(context.Context, Dto.UserLoginData) DomainLevel.LoginApplicationOutComingData
}

type RegisterApplication interface {
	RegisterApplication(context.Context, *Dto.UserDataRegister) DomainLevel.RegisterApplicationOutComingData
}
type UploadApplication interface {
	UploadApplication(UploadIncomeData) (string, error)
}
type UploaderEncryptApplication interface {
	UploadEncryptApplication(UploadEncryptIncomeData) (string, error)
}
