package Application

import (
	"Kaban/internal/InfrastructureLayer/Crypto"
	"Kaban/internal/InfrastructureLayer/FileTransferring/s3Repo"
	"Kaban/internal/InfrastructureLayer/FileTransferring/testFileTransfering"
	"Kaban/internal/InfrastructureLayer/Parsers"
	"Kaban/internal/InfrastructureLayer/RedisInteration"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"net/http/httptest"
	"testing"
)

const LongText = `ssssssssssssjh  j`

func TestCryptoFunctions(t *testing.T) {
	encr := Crypto.NewAesCtr()
	uploadEncrypt := UploadEncrypt{
		NewUploadEncryptCrypto: NewUploadEncryptCrypto{
			Encrypt: encr,
		},
	}
	downloadDecrypt := DownloadEncrypt{
		NewDownloadEncryptCrypto: NewDownloadEncryptCrypto{
			Aes: encr,
		},
	}
	firstUploading, _ := testFileTransfering.NewStartUploading().Make(nil)
	download, _ := testFileTransfering.NewTestDownload().Make(nil)
	uploadingEncr, _ := testFileTransfering.NewUploadingDecrypted().Make(nil)
	allData := make([]byte, encr.GetRequiredOverheadSize()+32)
	if _, err := io.ReadFull(rand.Reader, allData); err != nil {
		panic(err)
	}
	err := uploadEncrypt.encryptFile(bytes.NewReader([]byte(LongText)), firstUploading, allData)
	if err != nil {
		panic(err)
	}
	err = downloadDecrypt.decryptFile(allData[encr.GetRequiredOverheadSize():], download, uploadingEncr)
	if err != nil {
		panic(err)
	}
}

func TestApplicationLayer(t *testing.T) {
	encr := Crypto.NewAesGcmEncryption()
	rsa := Crypto.NewRsaEncryption()
	firstUploading := testFileTransfering.NewStartUploading()
	download := testFileTransfering.NewTestDownload()
	uploadingEncr := testFileTransfering.NewUploadingDecrypted()
	keys := RepoEncrypterKeys.GetKeysTest()
	redisTestIo := RedisInteration.NewRedisTestIO()
	uploadEncrypt := UploadEncrypt{
		NewUploadEncryptDataManage: NewUploadEncryptDataManage{
			Encode: Parsers.GetNewParsing(),
		},
		NewUploadEncryptCrypto: NewUploadEncryptCrypto{
			Generate:      Crypto.GetNewGenerating(),
			Encrypt:       encr,
			Rsa:           rsa,
			EncrypterKeys: keys,
		},
		NewUploadEncryptDelivery: NewUploadEncryptDelivery{
			Uploader:     firstUploading,
			RedisWriter:  redisTestIo,
			RedisChecker: RedisInteration.NewRedisCheckTest(),
			RedisDeleter: RedisInteration.NewRedisDelete(),
			Deleter:      s3Repo.NewS3DeleteTest(),
		},
	}
	downloadDecrypt := DownloadEncrypt{
		NewDownloadEncryptDelivery: NewDownloadEncryptDelivery{
			ReaderRedis:  redisTestIo,
			DeleterRedis: RedisInteration.NewRedisDeleteTest(),
		},
		NewDownloadEncryptCrypto: NewDownloadEncryptCrypto{
			Aes:           encr,
			Rsa:           rsa,
			EncrypterKeys: keys,
		},
		NewDownloadEncryptFileControl: NewDownloadEncryptFileControl{
			Uploader:   uploadingEncr,
			Downloader: download,
			Deleter:    s3Repo.NewS3DeleteTest(),
		},
		decoder: Parsers.GetNewParsing(),
	}

	_, err := uploadEncrypt.UploadEncrypt(IncomeData{
		File: bytes.NewReader([]byte(LongText)),
		Name: "da",
		Ctx:  context.Background(),
		Size: int64(len([]byte(LongText))),
	})

	if err != nil {
		t.Errorf("Got error %v\n", err)
		return
	}

	err = downloadDecrypt.DownloadEncrypt(NewDownloadEncryptIncomingData{
		dst:           httptest.NewRecorder(),
		Ctx:           context.Background(),
		EncryptedName: "",
	})
	if err != nil {
		t.Errorf("Got error %v\n", err)
		return
	}
}
