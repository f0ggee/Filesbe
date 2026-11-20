package testFunctions

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Drivers/Crypto"
	"Kaban/internal/Drivers/FileTransferring/testFileTransfering"
	"Kaban/internal/Drivers/Parsers"
	"Kaban/internal/Drivers/RedisInteration"
	"Kaban/internal/Drivers/RepoEncrypterKeys"
	"Kaban/internal/Service/Application"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

func TestEncryptUploading() Application.UploaderEncryptApplication {
	rsaK, _ := rsa.GenerateKey(rand.Reader, 2048)
	DomainLevel.GetCryptoKeys = RepoEncrypterKeys.NewTestEncrypterKeys(x509.MarshalPKCS1PrivateKey(rsaK))

	var fl = testFileTransfering.GetTestFlow()
	upload := Application.UploadEncrypt{
		NewUploadEncryptDataManage: Application.NewUploadEncryptDataManage{
			Encode: Parsers.GetNewParsing(),
		},
		NewUploadEncryptCrypto: Application.NewUploadEncryptCrypto{
			Generate: Crypto.GetNewGenerating(),
			Aes:      Crypto.NewAesGcmEncryption(),
			Rsa:      Crypto.NewRsaEncryption(),
		},
		NewUploadEncryptDelivery: Application.NewUploadEncryptDelivery{
			Uploader:     fl.U,
			RedisWriter:  RedisInteration.NewRedisWriteTest(),
			RedisChecker: RedisInteration.NewRedisCheckTest(),
			RedisDeleter: RedisInteration.NewRedisDeleteTest(),
			Deleter:      fl.D_,
		},
	}
	return &upload
}
