package cmds

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"Kaban/internal/InfrastructureLayer/Crypto"
	"Kaban/internal/InfrastructureLayer/DatabaseControl"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadEncryptRepo"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadNoEncrypt"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoLoginRealizations"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoRegisterRepository"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoSessionHandle"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoUsersCheckAuth"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepofileUploaderEncryptRepo"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepofileUploaderNoEncryptRepo"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepourlBuilder"
	"Kaban/internal/InfrastructureLayer/FileControls"
	"Kaban/internal/InfrastructureLayer/GrpcManage"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"Kaban/internal/InfrastructureLayer/s3Repo"
	"Kaban/internal/Service/Application"
	"os"
	"sync"

	"github.com/awnumar/memguard"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
)

type S3Collector struct {
	Deleter    s3Repo.DeleterS3
	Uploader   s3Repo.S3Uploader
	S3Download s3Repo.DownloadingS3
}

func GetS3RealizationsCollector(cfg *s3.Client, OldS3Connect *session.Session) *S3Collector {
	s3Info := s3Repo.GetNewVariables(os.Getenv("Bucket"), cfg, OldS3Connect)
	S3Upload := s3Repo.GetNewUploading(*s3Info)
	S3Download := s3Repo.GetNewS3Download(*s3Info)
	S3Deleter := s3Repo.GetNewDeleterS3(*s3Info)
	return &S3Collector{
		Deleter:    S3Deleter,
		Uploader:   S3Upload,
		S3Download: S3Download,
	}
}

type FileControlCollector struct {
	Transfer     FileControls.Transfer
	FileSettings FileControls.FileSettings
}

func GetFileControlCollector() *FileControlCollector {
	return &FileControlCollector{
		Transfer:     *FileControls.GetNewTransfer(),
		FileSettings: *FileControls.GetNewFileSettings(),
	}
}

func GetNewEncrypterKeysCollector(oldKey *memguard.LockedBuffer, newKey *memguard.LockedBuffer) *Application.EncrypterKeys {
	K := RepoEncrypterKeys.GetNewKeys(oldKey, newKey)
	return Application.GetEncrypterKeys(*K)
}

type CollectorCrypto struct {
	Validate DomainLevel.CryptoValidating
	Encrypt  DomainLevel.Encryption
	Decrypt  DomainLevel.Decryption
	Generate DomainLevel.CryptoGenerating
	Keys     DomainLevel.NewSetKeys
}

type NewCryptoCollectorInput struct {
	OurPrivateKey   []byte
	MasterPublicKet []byte
	Decode          RepoParsers.Decode
}

func GetNewCryptoCollector(d NewCryptoCollectorInput) *CollectorCrypto {

	Validation := Crypto.GetNeValidating()
	Encrypter := Crypto.GetNewEncrypter()
	Decrypter := Crypto.GetNewDecryption(d.Decode)
	Generate := Crypto.GetNewGenerating()
	Keys := DomainLevel.GetNewSetKeys(d.OurPrivateKey, d.MasterPublicKet)
	return &CollectorCrypto{
		Validate: Validation,
		Encrypt:  Encrypter,
		Decrypt:  Decrypter,
		Generate: Generate,
		Keys:     *Keys,
	}
}

type CollectorAuthTokensManage struct {
	Creating AuthTokensManage.CreatingTokens
	Validate AuthTokensManage.NewAuthChecker
}

func GetAuthTokensCollector(key []byte) *CollectorAuthTokensManage {
	creating := AuthTokensManage.GetNewCreatingTokens()
	Validate := AuthTokensManage.GetNNewAuthChecker(*creating, key)
	return &CollectorAuthTokensManage{
		Creating: *creating,
		Validate: *Validate,
	}
}

type CollectorDatabaseManage struct {
	Checker DatabaseControl.CheckerDb
	Reader  DatabaseControl.Read
	Writer  DatabaseControl.Writer
}

func GetDatabaseManageCollector(Db *pgxpool.Pool) CollectorDatabaseManage {
	Checker := DatabaseControl.GetNewCheckerDb(Db)
	Reader := DatabaseControl.GetNewRead(Db)
	Writer := DatabaseControl.GetNewWriter(Db)

	return CollectorDatabaseManage{
		Checker: *Checker,
		Reader:  *Reader,
		Writer:  *Writer,
	}
}

type DownloadEncryptCollector struct {
	Answ RepoDownloadEncryptRepo.NewFileDownloadEncrypt
}
type DownloadCollector struct {
	Answ RepoDownloadNoEncrypt.NewRepoDownloadNoEncrypt
}
type UploaderEncrypterCollector struct {
	Answ RepofileUploaderEncryptRepo.SetNewUploadingRepo
}
type UploaderCollector struct {
	Answ RepofileUploaderNoEncryptRepo.NewUploaderNoEncrypt
}
type LoginCollector struct {
	Answ RepoLoginRealizations.LoginAnswers
}
type RegisterCollector struct {
	Answ RepoRegisterRepository.NewRegister
}
type SessionCollector struct {
	Sess RepoSessionHandle.NewSessionConnect
}
type UrlBuilderCollector struct {
	Answ RepourlBuilder.NewUrlBuilder
}
type UserCheckCollector struct {
	Answ RepoUsersCheckAuth.SetUsersChecker
}
type DeliverPackagesCollector struct {
	DownloadEncryptCollector
	DownloadCollector
	UploaderEncrypterCollector
	UploaderCollector
	LoginCollector
	RegisterCollector
	SessionCollector
	UrlBuilderCollector
	UserCheckCollector
}

func GetDeliverPackagesCollector(Store *sessions.CookieStore) *DeliverPackagesCollector {

	return &DeliverPackagesCollector{
		DownloadEncryptCollector: DownloadEncryptCollector{
			Answ: *RepoDownloadEncryptRepo.GetNewFileDownloadEncrypt(),
		},
		DownloadCollector: DownloadCollector{
			Answ: *RepoDownloadNoEncrypt.GetNewNewRepoDownloadNoEncrypt(),
		},
		UploaderEncrypterCollector: UploaderEncrypterCollector{
			Answ: *RepofileUploaderEncryptRepo.GetNewSetNewUploadingRepo(),
		},
		UploaderCollector: UploaderCollector{
			Answ: *RepofileUploaderNoEncryptRepo.GetNewUploaderNoEncrypt(),
		},
		LoginCollector: LoginCollector{
			Answ: *RepoLoginRealizations.GetNewLoginAnswers(),
		},
		RegisterCollector: RegisterCollector{
			Answ: *RepoRegisterRepository.GetNewRegisterController(),
		},
		SessionCollector: SessionCollector{
			Sess: *RepoSessionHandle.GetNewSessionConnect(Store, &sync.RWMutex{}),
		},
		UrlBuilderCollector: UrlBuilderCollector{
			Answ: *RepourlBuilder.GetNewUrlBuilder(),
		},
		UserCheckCollector: UserCheckCollector{
			Answ: *RepoUsersCheckAuth.GetNewUsersChecker(),
		},
	}
}

type FileControlsCollector struct {
	Settings    FileControls.FileSettings
	Transfering FileControls.Transfer
}

func GetFileControlsCollector() *FileControlsCollector {
	return &FileControlsCollector{
		Settings:    *FileControls.GetNewFileSettings(),
		Transfering: *FileControls.GetNewTransfer(),
	}
}

type GrpcCollector struct {
	Sender   GrpcManage.NewSenderRequests
	Checking GrpcManage.HandlerGrpcRequest
}

func GetGrpcCollector() *GrpcCollector {

}
