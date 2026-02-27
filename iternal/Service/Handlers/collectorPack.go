package Handlers

import (
	"Kaban/iternal/DomainLevel"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type HandlerPack struct {
	Tokens    DomainLevel.ManageTokens
	Database  DomainLevel.UserServer
	TokenImpl DomainLevel.ManageTokensImpl
	RedisConn DomainLevel.RedisInteration
	S3Conn    DomainLevel.S3Interation
	FileInfo  DomainLevel.FileInfo
	Choose    DomainLevel.EncryptionKeyInteration

	S3Connect *s3.Client
}

type HandlerPackCollect struct {
	S HandlerPack
}

func CollectorPack(S HandlerPack) *HandlerPackCollect {
	return &HandlerPackCollect{S: S}
}
