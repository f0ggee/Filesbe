package Handlers

import (
	"Kaban/iternal/DomainLevel"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type HandlerPack struct {
	Tokens               DomainLevel.ManageAuthTokens
	Database             DomainLevel.UserServer
	TokenImpl            DomainLevel.ManageAuthTokensImpl
	RedisConn            DomainLevel.RedisInteraction
	S3Conn               DomainLevel.S3Handle
	FileInfo             DomainLevel.FileInfoManipulation
	Choose               DomainLevel.KeyInteraction
	GrpcDataMange        DomainLevel.GrpcDataManage
	Encryption           DomainLevel.EncryptionKey
	ConverterKey         DomainLevel.ConverterKey
	FileDataManipulation DomainLevel.FileDataManipulation
	S3Connect            *s3.Client
	GrpcConn             DomainLevel.GrpcInteraction
}

type HandlerPackCollect struct {
	S HandlerPack
}

func CollectorPack(S HandlerPack) *HandlerPackCollect {
	return &HandlerPackCollect{S: S}
}
