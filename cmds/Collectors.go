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
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoUsersCheckAuth"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepofileUploaderEncryptRepo"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepofileUploaderNoEncryptRepo"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepourlBuilder"
	"Kaban/internal/InfrastructureLayer/FileControls"
	"Kaban/internal/InfrastructureLayer/GrpcManage"
	"Kaban/internal/InfrastructureLayer/RedisInteration"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"Kaban/internal/InfrastructureLayer/RepoSessionHandle"
	"Kaban/internal/InfrastructureLayer/s3Repo"
	"os"
	"sync"

	"github.com/awnumar/memguard"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	ses "github.com/aws/aws-sdk-go/aws/session"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type S3Collector struct {
	Deleter    s3Repo.DeleterS3
	Uploader   s3Repo.S3Uploader
	S3Download s3Repo.DownloadingS3
}

func GetS3Collector(cfg *s3.Client, OldS3Connect *ses.Session) *S3Collector {
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

type CollectorCrypto struct {
	Validate DomainLevel.CryptoValidating
	Encrypt  DomainLevel.Encryption
	Decrypt  DomainLevel.Decryption
	Generate DomainLevel.CryptoGenerating
}

type NewCryptoCollectorInput struct {
	Decode RepoParsers.Decode
}

func GetNewCryptoCollector(d NewCryptoCollectorInput) *CollectorCrypto {

	Validation := Crypto.GetNeValidating()
	Encrypter := Crypto.GetNewEncrypter()
	Decrypter := Crypto.GetNewDecryption(d.Decode)
	Generate := Crypto.GetNewGenerating()
	return &CollectorCrypto{
		Validate: Validation,
		Encrypt:  Encrypter,
		Decrypt:  Decrypter,
		Generate: Generate,
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
		UrlBuilderCollector: UrlBuilderCollector{
			Answ: *RepourlBuilder.GetNewUrlBuilder(),
		},
		UserCheckCollector: UserCheckCollector{
			Answ: *RepoUsersCheckAuth.GetNewUsersChecker(),
		},
	}
}

type SessionCollector struct {
	Session RepoSessionHandle.NewSessionConnect
}

func GetSessionCollector(activity *sessions.CookieStore) *SessionCollector {
	return &SessionCollector{Session: *RepoSessionHandle.GetNewSessionConnect(activity, nil)}
}

type GrpcCollector struct {
	Sender   GrpcManage.NewSenderRequests
	Checking GrpcManage.HandlerGrpcRequest
}

func GetGrpcCollector(CryptoEncrypt DomainLevel.Encryption, CryptoDecrypt DomainLevel.Decryption, Parse RepoParsers.Decode, CryptoValidate DomainLevel.CryptoValidating, Keys RepoEncrypterKeys.Keys, ServerKeys DomainLevel.NewServerKeys) *GrpcCollector {

	return &GrpcCollector{
		Sender: *GrpcManage.GetNewSenderRequests(),
		Checking: *GrpcManage.GetNewHandlerGrpcRequest(GrpcManage.NewValidating{
			CryptoValidate: CryptoValidate,
		}, GrpcManage.NewKeys{
			Keys:       Keys,
			ServerKeys: ServerKeys,
		}, GrpcManage.NewDecrypt{
			CryptoDecrypt: CryptoDecrypt,
		}, GrpcManage.NewEncrypt{
			CryptoEncrypt: CryptoEncrypt,
		}, GrpcManage.NewParser{
			Parse: Parse,
		}),
	}
}

type RedisCollector struct {
	Delete RedisInteration.DeleterRedis
	Write  RedisInteration.Writing
	Read   RedisInteration.RedisReader
	Check  RedisInteration.ValidationRedis
}

func GetRedisCollector(Re *redis.Client) *RedisCollector {
	return &RedisCollector{
		Delete: *RedisInteration.GetNewDeleterRedis(Re),
		Write:  *RedisInteration.GetNewWriting(Re),
		Read:   *RedisInteration.GetNewRedisReader(Re),
		Check:  *RedisInteration.GetNewValidationRedis(Re),
	}
}

func GetNewServerKeysCollector(key1 []byte, key2 []byte) *DomainLevel.NewServerKeys {
	return DomainLevel.GetNewSetKeys(key1, key2)
}

func GetEncrypterKeysCollector(Key1 *memguard.LockedBuffer, Key2 *memguard.LockedBuffer) *RepoEncrypterKeys.Keys {
	return RepoEncrypterKeys.GetNewKeys(Key1, Key2)
}

type RepoParsersCollector struct {
	Decode RepoParsers.Decode
	Encode RepoParsers.Encode
}

func GetRepoParsersCollector() *RepoParsersCollector {
	return &RepoParsersCollector{
		Decode: RepoParsers.GetNewParsing(),
		Encode: RepoParsers.GetNewParsing(),
	}
}
