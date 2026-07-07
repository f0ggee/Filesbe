package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadEncryptRepo"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadNoEncrypt"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"Kaban/internal/InfrastructureLayer/s3Repo"
)

type s3Controlling struct {
	Deleter    s3Repo.DeleterS3
	Uploader   s3Repo.S3Uploader
	S3Download s3Repo.DownloadingS3
}
type FileDownload struct {
	Download RepoDownloadNoEncrypt.NewDownloadFile
}
type fileDownloadEncrypt struct {
	EncryptDownload RepoDownloadEncryptRepo.NewEncryptDownloadFile
}

type encrypterKeys struct {
	GetKeys RepoEncrypterKeys.Keys
}
type getCrypto struct {
	Validate DomainLevel.CryptoValidating
	Encrypt  DomainLevel.Encryption
	Decrypt  DomainLevel.Decryption
	Generate DomainLevel.CryptoGenerating
	Keys     DomainLevel.NewSetKeys
}

type getFileManager struct {
	FileManaging DomainLevel.SetFileSettings
}

type getControlKeys struct {
	Keys DomainLevel.NewSetKeys
}
type authTokens struct {
	GeneratingToken AuthTokensManage.Generator
	Checking        AuthTokensManage.AuthCheck
}
type HandlerGrpc struct {
	GrpcSendingRequest DomainLevel.SendRequestGrpc
	ProcessingRequests DomainLevel.HandlingRequests
}
type databaseControlling struct {
	Writer  DomainLevel.WriteDb
	Reader  DomainLevel.ReadDb
	Checker DomainLevel.CheckingDb
}
type redisControlling struct {
	Writer       DomainLevel.WritingRedis
	Reader       DomainLevel.ReadingRedis
	Deleter      DomainLevel.DeleterRedis
	CheckerRedis DomainLevel.RedisChecker
}

type KeysControlling struct {
	ControllerKey DomainLevel.NewSetKeys
}

type parser struct {
	Decode RepoParsers.Decode
	Encode RepoParsers.Encode
}
