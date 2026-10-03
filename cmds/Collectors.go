//A collector layer is a collection of functions that build
//necessary modules. This layer is an entry point of a program

package cmds

import (
	"Kaban/internal/InfrastructureLayer/Crypto"
	"Kaban/internal/InfrastructureLayer/DatabaseControl"
	"Kaban/internal/InfrastructureLayer/FileTransferring/TransferringHttpRepo"
	"Kaban/internal/InfrastructureLayer/FileTransferring/s3Repo"
	grcRec "Kaban/internal/InfrastructureLayer/Grpc/GrpcRequests"
	"Kaban/internal/InfrastructureLayer/Parsers"
	"Kaban/internal/InfrastructureLayer/RedisInteration"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"Kaban/internal/InfrastructureLayer/Tokens"
)

type CryptoCollector struct {
	RsaRealization    Crypto.RsaEncryption
	AesRealization    Crypto.AesEncryption
	AesCtrRealization Crypto.AesCtrEncryption
}

func GetCryptoCollector() CryptoCollector {
	return CryptoCollector{
		RsaRealization:    Crypto.NewRsaEncryption(),
		AesRealization:    Crypto.NewAesEncryption(),
		AesCtrRealization: Crypto.NewAesCtr(),
	}
}

type DatabaseCollector struct {
	Check *DatabaseControl.CheckerDb
	Read  *DatabaseControl.Read
	Write *DatabaseControl.Writer
}

func GetDatabaseCollector(conn DatabaseControl.DatabaseConn) DatabaseCollector {
	return DatabaseCollector{
		Check: DatabaseControl.NewCheckerDb(conn),
		Read:  DatabaseControl.NewRead(conn),
		Write: DatabaseControl.NewWriter(conn),
	}
}

type S3FileTransferringCollector struct {
	Upload   s3Repo.S3Upload
	Download s3Repo.S3Download
	Delete   s3Repo.S3Delete
}
type HttpFileTransferringCollector struct {
	Upload TransferringHttpRepo.HttpUploader
}
type FileTransferringCollector struct {
	S3FileTransferringCollector
	HttpFileTransferringCollector
}

func GetFileTransferring() FileTransferringCollector {
	return FileTransferringCollector{
		S3FileTransferringCollector: S3FileTransferringCollector{
			Upload:   s3Repo.NewS3Upload(),
			Download: s3Repo.NewS3Downloader(),
			Delete:   s3Repo.NewS3Deleter(),
		},
		HttpFileTransferringCollector: HttpFileTransferringCollector{
			Upload: TransferringHttpRepo.NewHttpUploader(),
		},
	}
}

type GrpcCollector struct {
	Reqs grcRec.SetNewKeyRequest
}

func GetGrpcCollector() GrpcCollector {
	return GrpcCollector{
		Reqs: grcRec.GetNewSetNewKeyRequest(),
	}
}

type ParsersCollector struct {
	Encode Parsers.Parsing
}

func NewParsersCollector() ParsersCollector {
	return ParsersCollector{
		Encode: Parsers.GetNewParsing(),
	}
}

type Redis struct {
	Read   RedisInteration.RedisRead
	Write  RedisInteration.RedisWrite
	Check  RedisInteration.RedisCheck
	Delete RedisInteration.RedisDelete
}

func NewRedis() *Redis {
	return &Redis{
		Read:   RedisInteration.NewRedisRead(),
		Write:  RedisInteration.NewRedisWrite(),
		Check:  RedisInteration.NewRedisCheck(),
		Delete: RedisInteration.NewRedisDelete(),
	}
}

type SessionKeys struct {
	Keys RepoEncrypterKeys.Keys
}

func NewSessionKeys() SessionKeys {
	return SessionKeys{
		Keys: RepoEncrypterKeys.GetNewKeys(),
	}
}

type TokensAuth struct {
	Jwt Tokens.JwtToken
	Rf  Tokens.RfToken
}

func NewTokensAuth() TokensAuth {
	return TokensAuth{
		Jwt: Tokens.GetNewJwtToken(),
		Rf:  Tokens.GetNewRfToken(),
	}
}
