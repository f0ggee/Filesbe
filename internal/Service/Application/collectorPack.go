package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadEncryptRepo"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadNoEncrypt"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoParsers"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
)

type S3Controlling struct {
	Deleter    DomainLevel.DeleterS3
	Uploader   DomainLevel.S3Uploader
	S3Download DomainLevel.DownloadingS3
}
type FileDownload struct {
	Download RepoDownloadNoEncrypt.NewDownloadFile
}
type FileDownloadEncrypt struct {
	EncryptDownload RepoDownloadEncryptRepo.NewEncryptDownloadFile
}

type EncrypterKeys struct {
	GetKeys RepoEncrypterKeys.Keys
}
type GetCrypto struct {
	Validate DomainLevel.CryptoValidating
	Encrypt  DomainLevel.Encryption
	Decrypt  DomainLevel.Decryption
	Generate DomainLevel.CryptoGenerating
	Keys     DomainLevel.CryptoKey
}

type GetFileManager struct {
	FileInfo     DomainLevel.HandleFileInfo
	FileManaging DomainLevel.HandleFile
}

type GetControlKeys struct {
	Keys DomainLevel.NewSetKeys
}
type AuthTokens struct {
	Manage          DomainLevel.ManageTokens
	GeneratingToken DomainLevel.Generator
	Checking        DomainLevel.CheckingAuthTokens
}
type HandlerGrpc struct {
	GrpcSendingRequest DomainLevel.SendRequestGrpc
	ProcessingRequests DomainLevel.HandlingRequests
	GrpcTest           DomainLevel.GrpcTest
}
type DatabaseControlling struct {
	Writer  DomainLevel.WriteDb
	Reader  DomainLevel.ReadDb
	Checker DomainLevel.CheckingDb
}
type RedisControlling struct {
	Writer       DomainLevel.WritingRedis
	Reader       DomainLevel.ReadingRedis
	Deleter      DomainLevel.DeleterRedis
	CheckerRedis DomainLevel.RedisChecker
}

type KeysControlling struct {
	ControllerKey DomainLevel.NewSetKeys
}
type Converter struct {
	Converting DomainLevel.DataConvert
}

type Parser struct {
	Decode RepoParsers.Decode
	Encode RepoParsers.Encode
}
