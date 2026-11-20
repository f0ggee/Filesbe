package testFunctions

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Drivers/Crypto"
	"Kaban/internal/Drivers/FileTransferring/TransferringHttpRepo"
	"Kaban/internal/Drivers/FileTransferring/s3Repo"
	"Kaban/internal/Drivers/FileTransferring/testFileTransfering"
	"Kaban/internal/Drivers/Parsers"
	"Kaban/internal/Drivers/RedisInteration"
	"Kaban/internal/Drivers/RepoEncrypterKeys"
	"Kaban/internal/Service/Application"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io"
)

func DownloadEncryptTest(inputData []byte) Application.DownloadEncryptApplication {
	data := make([]byte, 32)
	io.ReadFull(rand.Reader, data)
	aes, _ := Crypto.NewAesGcmEncryption().MakeCrypto(data, nil)
	cipherText, _ := aes.Encrypt(inputData)
	rsaK, _ := rsa.GenerateKey(rand.Reader, 2048)

	da, err := Parsers.GetNewParsing().Encode(DomainLevel.FileLabelsBytes{
		FileName: []byte("Hello!"),
		AesKey:   data,
	})
	if err != nil {
		panic(err)
	}

	DomainLevel.GetCryptoKeys = RepoEncrypterKeys.NewTestEncrypterKeys(x509.MarshalPKCS1PrivateKey(rsaK))

	rsaC, _ := Crypto.NewRsaEncryption().MakeCrypto(x509.MarshalPKCS1PrivateKey(rsaK), nil)
	labels, _ := rsaC.Encrypt(da)

	flow := testFileTransfering.GetTestFlow()
	_u, _ := flow.U.Make(nil)
	_u.Upload(bytes.NewReader(cipherText))

	return &Application.DownloadEncrypt{
		NewDownloadEncryptDelivery: Application.NewDownloadEncryptDelivery{
			ReaderRedis:  RedisInteration.NewRedisTestDownloadEncrypt(labels),
			DeleterRedis: RedisInteration.NewRedisDelete(),
		},
		NewDownloadEncryptCrypto: Application.NewDownloadEncryptCrypto{
			Aes: Crypto.NewAesGcmEncryption(),
			Rsa: Crypto.NewRsaEncryption(),
		},
		NewDownloadEncryptFileControl: Application.NewDownloadEncryptFileControl{
			Uploader:   TransferringHttpRepo.NewHttpUploader(),
			Downloader: flow.D,
			Deleter:    s3Repo.NewS3DeleteTest(),
		},
		Decoder: Parsers.GetNewParsing(),
	}
}

// DownloadEncryptBadRsa test bad RSA keys
func DownloadEncryptBadRsa(inputData []byte) Application.DownloadEncryptApplication {
	data := make([]byte, 32)
	io.ReadFull(rand.Reader, data)
	aes, _ := Crypto.NewAesGcmEncryption().MakeCrypto(data, nil)
	cipherText, _ := aes.Encrypt(inputData)
	rsaK, _ := rsa.GenerateKey(rand.Reader, 2048)

	da, err := Parsers.GetNewParsing().Encode(DomainLevel.FileLabelsBytes{
		FileName: []byte("Hello!"),
		AesKey:   data,
	})
	if err != nil {
		panic(err)
	}
	badRsa, _ := rsa.GenerateKey(rand.Reader, 2048)
	DomainLevel.GetCryptoKeys = RepoEncrypterKeys.NewTestEncrypterKeys(x509.MarshalPKCS1PrivateKey(badRsa))

	rsaC, _ := Crypto.NewRsaEncryption().MakeCrypto(x509.MarshalPKCS1PrivateKey(rsaK), nil)
	labels, _ := rsaC.Encrypt(da)

	flow := testFileTransfering.GetTestFlow()
	_u, _ := flow.U.Make(nil)
	_u.Upload(bytes.NewReader(cipherText))

	return &Application.DownloadEncrypt{
		NewDownloadEncryptDelivery: Application.NewDownloadEncryptDelivery{
			ReaderRedis:  RedisInteration.NewRedisTestDownloadEncrypt(labels),
			DeleterRedis: RedisInteration.NewRedisDelete(),
		},
		NewDownloadEncryptCrypto: Application.NewDownloadEncryptCrypto{
			Aes: Crypto.NewAesGcmEncryption(),
			Rsa: Crypto.NewRsaEncryption(),
		},
		NewDownloadEncryptFileControl: Application.NewDownloadEncryptFileControl{
			Uploader:   TransferringHttpRepo.NewHttpUploader(),
			Downloader: flow.D,
			Deleter:    s3Repo.NewS3DeleteTest(),
		},
		Decoder: Parsers.GetNewParsing(),
	}

}
func DownloadEncryptBadAesKey(inputData []byte) Application.DownloadEncryptApplication {
	trueKey := make([]byte, 32)
	badKey := make([]byte, 32)
	io.ReadFull(rand.Reader, trueKey)
	io.ReadFull(rand.Reader, badKey)

	aes, _ := Crypto.NewAesGcmEncryption().MakeCrypto(trueKey, nil)
	cipherText, _ := aes.Encrypt(inputData)
	rsaK, _ := rsa.GenerateKey(rand.Reader, 2048)

	da, err := Parsers.GetNewParsing().Encode(DomainLevel.FileLabelsBytes{
		FileName: []byte("Hello!"),
		AesKey:   badKey,
	})
	if err != nil {
		panic(err)
	}

	DomainLevel.GetCryptoKeys = RepoEncrypterKeys.NewTestEncrypterKeys(x509.MarshalPKCS1PrivateKey(rsaK))

	rsaC, _ := Crypto.NewRsaEncryption().MakeCrypto(x509.MarshalPKCS1PrivateKey(rsaK), nil)
	labels, _ := rsaC.Encrypt(da)

	flow := testFileTransfering.GetTestFlow()
	_u, _ := flow.U.Make(nil)
	_u.Upload(bytes.NewReader(cipherText))

	return &Application.DownloadEncrypt{
		NewDownloadEncryptDelivery: Application.NewDownloadEncryptDelivery{
			ReaderRedis:  RedisInteration.NewRedisTestDownloadEncrypt(labels),
			DeleterRedis: RedisInteration.NewRedisDelete(),
		},
		NewDownloadEncryptCrypto: Application.NewDownloadEncryptCrypto{
			Aes: Crypto.NewAesGcmEncryption(),
			Rsa: Crypto.NewRsaEncryption(),
		},
		NewDownloadEncryptFileControl: Application.NewDownloadEncryptFileControl{
			Uploader:   TransferringHttpRepo.NewHttpUploader(),
			Downloader: flow.D,
			Deleter:    s3Repo.NewS3DeleteTest(),
		},
		Decoder: Parsers.GetNewParsing(),
	}
}
