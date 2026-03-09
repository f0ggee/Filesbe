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
	FileInfo  DomainLevel.FileInfoInteraction
	Choose    DomainLevel.KeyInteraction

	Encryption           DomainLevel.EncryptionKey
	FileDataManipulation DomainLevel.FileInfoDataManipulation
	S3Connect            *s3.Client
}

type HandlerPackCollect struct {
	S HandlerPack
}

func CollectorPack(S HandlerPack) *HandlerPackCollect {
	return &HandlerPackCollect{S: S}
}
