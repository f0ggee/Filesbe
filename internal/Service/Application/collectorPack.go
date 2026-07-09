package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"Kaban/internal/InfrastructureLayer/FileControls"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"Kaban/internal/InfrastructureLayer/s3Repo"
)

type S3Controlling struct {
	Deleter    s3Repo.DeleterS3
	Uploader   s3Repo.S3Uploader
	S3Download s3Repo.DownloadingS3
}

func GetNewS3Controlling(deleter s3Repo.DeleterS3, uploader s3Repo.S3Uploader, s3Download s3Repo.DownloadingS3) *S3Controlling {
	return &S3Controlling{Deleter: deleter, Uploader: uploader, S3Download: s3Download}
}

type Transfers struct {
	Transfer FileControls.Transfer
}

func GetNewTransfers(transfer FileControls.Transfer) *Transfers {
	return &Transfers{Transfer: transfer}
}

type EncrypterKeys struct {
	GetKeys RepoEncrypterKeys.Keys
}

func GetEncrypterKeys(getKeys RepoEncrypterKeys.Keys) *EncrypterKeys {
	return &EncrypterKeys{GetKeys: getKeys}
}

type GetCrypto struct {
	Validate DomainLevel.CryptoValidating
	Encrypt  DomainLevel.Encryption
	Decrypt  DomainLevel.Decryption
	Generate DomainLevel.CryptoGenerating
	Keys     DomainLevel.NewSetKeys
}

func GetNewCrypto(validate DomainLevel.CryptoValidating, encrypt DomainLevel.Encryption, decrypt DomainLevel.Decryption, generate DomainLevel.CryptoGenerating, keys DomainLevel.NewSetKeys) *GetCrypto {
	return &GetCrypto{Validate: validate, Encrypt: encrypt, Decrypt: decrypt, Generate: generate, Keys: keys}
}

type NewFileManager struct {
	FileManaging DomainLevel.SetFileSettings
}

func GetNewFileManager(fileManaging DomainLevel.SetFileSettings) *NewFileManager {
	return &NewFileManager{FileManaging: fileManaging}
}

type GetControlKeys struct {
	Keys DomainLevel.NewSetKeys
}

func GetNewControlKeys(keys DomainLevel.NewSetKeys) *GetControlKeys {
	return &GetControlKeys{Keys: keys}
}

type AuthTokens struct {
	GeneratingToken AuthTokensManage.Generator
	Checking        AuthTokensManage.AuthCheck
}

func GetNewAuthTokens(generatingToken AuthTokensManage.Generator, checking AuthTokensManage.AuthCheck) *AuthTokens {
	return &AuthTokens{GeneratingToken: generatingToken, Checking: checking}
}

type HandlerGrpc struct {
	GrpcSendingRequest DomainLevel.SendRequestGrpc
	ProcessingRequests DomainLevel.HandlingRequests
}

func GetNewHandlerGrpc(grpcSendingRequest DomainLevel.SendRequestGrpc, processingRequests DomainLevel.HandlingRequests) *HandlerGrpc {
	return &HandlerGrpc{GrpcSendingRequest: grpcSendingRequest, ProcessingRequests: processingRequests}
}

type DatabaseControlling struct {
	Writer  DomainLevel.WriteDb
	Reader  DomainLevel.ReadDb
	Checker DomainLevel.CheckingDb
}

func GetDatabaseControlling(writer DomainLevel.WriteDb, reader DomainLevel.ReadDb, checker DomainLevel.CheckingDb) *DatabaseControlling {
	return &DatabaseControlling{Writer: writer, Reader: reader, Checker: checker}
}

type RedisControlling struct {
	Writer       DomainLevel.WritingRedis
	Reader       DomainLevel.ReadingRedis
	Deleter      DomainLevel.DeleterRedis
	CheckerRedis DomainLevel.RedisChecker
}

func GetRedisControlling(writer DomainLevel.WritingRedis, reader DomainLevel.ReadingRedis, deleter DomainLevel.DeleterRedis, checkerRedis DomainLevel.RedisChecker) *RedisControlling {
	return &RedisControlling{Writer: writer, Reader: reader, Deleter: deleter, CheckerRedis: checkerRedis}
}

type Parser struct {
	Decode RepoParsers.Decode
	Encode RepoParsers.Encode
}

func GetParser(decode RepoParsers.Decode, encode RepoParsers.Encode) *Parser {
	return &Parser{Decode: decode, Encode: encode}
}
