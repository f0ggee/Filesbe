package Handlers

import (
	"Kaban/internal/DomainLevel"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
)

type S3Controlling struct {
	Deleter      DomainLevel.DeleterS3
	S3Connect    *s3.Client
	S3OldConnect *session.Session
}

type HandlerPackCrypto struct {
	Validate DomainLevel.CryptoValidating
	Encrypt  DomainLevel.Encryption
	Decrypt  DomainLevel.Decryption
	Generate DomainLevel.CryptoGenerating
}

type HandlerFileManagerPack struct {
	FileInfo     DomainLevel.HandleFileInfo
	FileManaging DomainLevel.HandleFile
}

type HandlerPackAuthTokens struct {
	Manage          DomainLevel.ManageTokens
	GeneratingToken DomainLevel.Generator
	Checking        DomainLevel.CheckingAuthTokens
}
type HandlerGrpc struct {
	GrpcSendingRequest DomainLevel.SendingRequestGrpc
	ProcessingRequests DomainLevel.HandlingRequests
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
	ControllerKey DomainLevel.KeysManager
}
type Converter struct {
	Converting DomainLevel.DataConvert
}
type HandlerPackCollect struct {
	S3                  S3Controlling
	Crypto              HandlerPackCrypto
	FileInfo            HandlerFileManagerPack
	AuthTokens          HandlerPackAuthTokens
	DatabaseControlling DatabaseControlling
	RedisControlling    RedisControlling
	Grpc                HandlerGrpc
	Convert             Converter
	Keys                KeysControlling
}

func NewHandlerPackCollect(s3 S3Controlling, crypto HandlerPackCrypto, fileInfo HandlerFileManagerPack, authTokens HandlerPackAuthTokens, databaseControlling DatabaseControlling, redisControlling RedisControlling, grpc HandlerGrpc, convert Converter, keys KeysControlling) *HandlerPackCollect {
	return &HandlerPackCollect{S3: s3, Crypto: crypto, FileInfo: fileInfo, AuthTokens: authTokens, DatabaseControlling: databaseControlling, RedisControlling: redisControlling, Grpc: grpc, Convert: convert, Keys: keys}
}
