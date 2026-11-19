package testFunctions

import (
	"Kaban/internal/Drivers/Crypto"
	"Kaban/internal/Drivers/FileTransferring/testFileTransfering"
	"Kaban/internal/Drivers/Parsers"
	"Kaban/internal/Drivers/RedisInteration"
	"Kaban/internal/Service/Application"
)

func TestUpload() Application.UploadApplication {
	u := Application.NewUpload{
		NewFileUploaderCrypto: Application.NewFileUploaderCrypto{
			Generator: Crypto.GetNewGenerating(),
		},
		NewFileUploaderDataMange: Application.NewFileUploaderDataMange{
			Encode: Parsers.GetNewParsing(),
		},
		NewFileUploaderDelivery: Application.NewFileUploaderDelivery{
			Uploader:   testFileTransfering.GetNormalFlow().U,
			WriteRedis: RedisInteration.NewRedisWriteTest(),
		},
	}
	return &u
}
