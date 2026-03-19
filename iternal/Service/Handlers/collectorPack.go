package Handlers

import (
	"Kaban/iternal/DomainLevel"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type HandlerPack struct {
	Tokens               DomainLevel.ManageTokens
	Database             DomainLevel.ReadDb
	TokenImpl            DomainLevel.ManageAuthTokensImpl
	RedisConn            DomainLevel.DeleterRedis
	S3Conn               DomainLevel.S3Handle
	FileInfo             DomainLevel.FileInfoManipulation
	Choose               DomainLevel.KeyInteraction
	GrpcDataMange        DomainLevel.GrpcDataManage
	Encryption           DomainLevel.Encryption
	ConverterKey         DomainLevel.Converter
	FileDataManipulation DomainLevel.FileDataManipulation
	S3Connect            *s3.Client
	GrpcConn             DomainLevel.SendingRequestGrpc
	Checking             DomainLevel.PacketChecker
}

type HandlerPackCollect struct {
	S HandlerPack
}

func CollectorPack(S HandlerPack) *HandlerPackCollect {
	return &HandlerPackCollect{S: S}
}
