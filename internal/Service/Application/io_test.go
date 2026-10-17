package Application

import (
	"Kaban/internal/DomainLevel"
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

func TestIoCrypto(t *testing.T) {
	crypt := Crypto.NewAesCtr()
	fistUpload, _ := testFileTransfering.GetFirstUploading().Make(nil)
	download := testFileTransfering.GetTestDownloader()
	secondUpload := testFileTransfering.GetUploaderSecond()

	key := make([]byte, crypt.GetRequiredOverheadSize()+32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}

	type IncomeUploadEncrypt struct {
		src io.Reader
		dst DomainLevel.Upload
		key []byte
	}
	type IncomeDownloadEncrypt struct {
		key    []byte
		o      DomainLevel.Download
		upload DomainLevel.Upload
	}
	tests := []struct {
		name          string
		uploadParam   IncomeUploadEncrypt
		downloadParam IncomeDownloadEncrypt
		uploadField   UploadEncrypt
		downloadField DownloadEncrypt
		wantErr       bool
	}{
		{
			name: "test1",
			uploadParam: IncomeUploadEncrypt{
				src: bytes.NewReader([]byte("It's example")),
				dst: fistUpload,
				key: key,
			},
			downloadParam: IncomeDownloadEncrypt{
				key:    key,
				o:      download,
				upload: &secondUpload,
			},
			uploadField: UploadEncrypt{
				NewUploadEncryptCrypto: NewUploadEncryptCrypto{
					Encrypt: crypt,
				},
			},
			downloadField: DownloadEncrypt{
				NewDownloadEncryptCrypto: NewDownloadEncryptCrypto{
					Aes: crypt,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			err := tt.uploadField.encryptFile(tt.uploadParam.src, tt.uploadParam.dst, tt.uploadParam.key)
			if err != nil && !tt.wantErr {
				t.Errorf("We wanted %v but got %v\n", tt.wantErr, err)
			}

			err = tt.downloadField.decryptFile(tt.downloadParam.key[crypt.GetRequiredOverheadSize():], tt.downloadParam.o, tt.downloadParam.upload)
			if err != nil && !tt.wantErr {
				t.Errorf("We wanted %v but got %v\n", tt.wantErr, err)
			}

		})
	}
}

func TestApplicationIOCrypto(t *testing.T) {

	crypt := Crypto.NewAesCtr()
	fistUpload := testFileTransfering.GetFirstUploading()
	download := testFileTransfering.GetTestDownloader()
	secondUpload := testFileTransfering.GetUploaderSecond()
	redisIOTest := RedisInteration.NewRedisTestIO()

	encrypterKeysTest := RepoEncrypterKeys.GetKeysTest()
	key := make([]byte, crypt.GetRequiredOverheadSize()+32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}

	type IncomeUploadEncrypt struct {
		src io.Reader
		dst DomainLevel.Upload
		key []byte
	}
	type IncomeDownloadEncrypt struct {
		key    []byte
		o      DomainLevel.Download
		upload DomainLevel.Upload
	}
	tests := []struct {
		name          string
		uploadParam   IncomeUploadEncrypt
		downloadParam IncomeDownloadEncrypt
		uploadField   UploadEncrypt
		downloadField DownloadEncrypt
		wantErr       bool
	}{
		{
			name: "test1",
			uploadParam: IncomeUploadEncrypt{
				src: bytes.NewReader([]byte("It's example")),
				dst: &fistUpload,
				key: key,
			},
			downloadParam: IncomeDownloadEncrypt{
				key:    key,
				o:      download,
				upload: &secondUpload,
			},
			uploadField: UploadEncrypt{
				NewUploadEncryptDataManage: NewUploadEncryptDataManage{
					Encode: Parsers.GetNewParsing(),
				},
				NewUploadEncryptCrypto: NewUploadEncryptCrypto{
					Generate:      Crypto.GetNewGenerating(),
					Encrypt:       crypt,
					Rsa:           Crypto.NewRsaEncryption(),
					EncrypterKeys: encrypterKeysTest,
				},
				NewUploadEncryptDelivery: NewUploadEncryptDelivery{
					Uploader:     fistUpload,
					RedisWriter:  redisIOTest,
					RedisChecker: RedisInteration.NewRedisCheckTest(),
					RedisDeleter: RedisInteration.NewRedisDeleteTest(),
					Deleter:      s3Repo.NewS3DeleteTest(),
				},
			},
			downloadField: DownloadEncrypt{
				NewDownloadEncryptDelivery: NewDownloadEncryptDelivery{
					ReaderRedis:  redisIOTest,
					DeleterRedis: RedisInteration.NewRedisDeleteTest(),
				},
				NewDownloadEncryptCrypto: NewDownloadEncryptCrypto{
					Aes:           crypt,
					Rsa:           Crypto.NewRsaEncryption(),
					EncrypterKeys: encrypterKeysTest,
				},
				NewDownloadEncryptFileControl: NewDownloadEncryptFileControl{
					Uploader:   &secondUpload,
					Downloader: download,
					Deleter:    s3Repo.NewS3DeleteTest(),
				},
				decoder: Parsers.GetNewParsing(),
			},
			wantErr: false,
		},
		{
			name: "test1",
			uploadParam: IncomeUploadEncrypt{
				src: bytes.NewReader([]byte("dkjnhssssssssss kkkkkdskdksdkskdksdkskdksdks kdjsj dkjsdj ksjd kjskdj ksjdkj skjdksdksjdkskdskdskdjaksjdklaksdaklsdlkasdklaslkfaklflkafklfklkmfvkvdfnkjfkdfjkdnvdjfngvjdkfjk")),
				dst: &fistUpload,
				key: key,
			},
			downloadParam: IncomeDownloadEncrypt{
				key:    key,
				o:      download,
				upload: &secondUpload,
			},
			uploadField: UploadEncrypt{
				NewUploadEncryptDataManage: NewUploadEncryptDataManage{
					Encode: Parsers.GetNewParsing(),
				},
				NewUploadEncryptCrypto: NewUploadEncryptCrypto{
					Generate:      Crypto.GetNewGenerating(),
					Encrypt:       crypt,
					Rsa:           Crypto.NewRsaEncryption(),
					EncrypterKeys: encrypterKeysTest,
				},
				NewUploadEncryptDelivery: NewUploadEncryptDelivery{
					Uploader:     fistUpload,
					RedisWriter:  redisIOTest,
					RedisChecker: RedisInteration.NewRedisCheckTest(),
					RedisDeleter: RedisInteration.NewRedisDeleteTest(),
					Deleter:      s3Repo.NewS3DeleteTest(),
				},
			},
			downloadField: DownloadEncrypt{
				NewDownloadEncryptDelivery: NewDownloadEncryptDelivery{
					ReaderRedis:  redisIOTest,
					DeleterRedis: RedisInteration.NewRedisDeleteTest(),
				},
				NewDownloadEncryptCrypto: NewDownloadEncryptCrypto{
					Aes:           crypt,
					Rsa:           Crypto.NewRsaEncryption(),
					EncrypterKeys: encrypterKeysTest,
				},
				NewDownloadEncryptFileControl: NewDownloadEncryptFileControl{
					Uploader:   &secondUpload,
					Downloader: download,
					Deleter:    s3Repo.NewS3DeleteTest(),
				},
				decoder: Parsers.GetNewParsing(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.uploadField.UploadEncrypt(UploadEncryptIncomeData{
				File: tt.uploadParam.src,
				Name: "",
				Ctx:  context.Background(),
				Size: 100,
			})
			if err != nil && !tt.wantErr {
				t.Fatalf("Got an error %v but wanted %v", err, tt.wantErr)
			}

			err = tt.downloadField.DownloadEncrypt(DownloadEncryptIncomeData{
				dst:           httptest.NewRecorder(),
				Ctx:           context.Background(),
				EncryptedName: "",
			})
		})
	}

}

func TestApplicationIO(t *testing.T) {
	fistUpload := testFileTransfering.GetFirstUploading()
	download := testFileTransfering.GetTestDownloader()
	secondUpload := testFileTransfering.GetUploaderSecond()
	redisIOTest := RedisInteration.NewRedisTestIO()
	tests := []struct {
		name          string
		Upload        NewUpload
		Download      NewDownload
		UploadParam   UploadIncomeData
		DownloadParam struct {
			name string
			ctx  context.Context
		}
	}{
		{
			name: "test1",
			Upload: NewUpload{
				NewFileUploaderCrypto: NewFileUploaderCrypto{
					Generator: Crypto.GetNewGenerating(),
				},
				NewFileUploaderDataMange: NewFileUploaderDataMange{
					Encode: Parsers.GetNewParsing(),
				},
				NewFileUploaderDelivery: NewFileUploaderDelivery{
					Uploader:   fistUpload,
					WriteRedis: redisIOTest,
				},
			},
			Download: NewDownload{
				DownloadFileControl: DownloadFileControl{
					Downloader: download,
					Uploader:   &secondUpload,
					Deleter:    testFileTransfering.NewTestDelete(),
					Decoder:    Parsers.GetNewParsing(),
				},
				Reader: redisIOTest,
				W:      httptest.NewRecorder(),
			},
			UploadParam: UploadIncomeData{
				File: bytes.NewReader([]byte("It's example")),
				Name: "test",
				Size: 100,
				Ctx:  context.Background(),
			},
			DownloadParam: struct {
				name string
				ctx  context.Context
			}{
				name: "test",
				ctx:  context.Background(),
			},
		},
		{
			name: "test1",
			Upload: NewUpload{
				NewFileUploaderCrypto: NewFileUploaderCrypto{
					Generator: Crypto.GetNewGenerating(),
				},
				NewFileUploaderDataMange: NewFileUploaderDataMange{
					Encode: Parsers.GetNewParsing(),
				},
				NewFileUploaderDelivery: NewFileUploaderDelivery{
					Uploader:   fistUpload,
					WriteRedis: redisIOTest,
				},
			},
			Download: NewDownload{
				DownloadFileControl: DownloadFileControl{
					Downloader: download,
					Uploader:   &secondUpload,
					Deleter:    testFileTransfering.NewTestDelete(),
					Decoder:    Parsers.GetNewParsing(),
				},
				Reader: redisIOTest,
				W:      httptest.NewRecorder(),
			},
			UploadParam: UploadIncomeData{
				File: bytes.NewReader([]byte("It's example")),
				Name: "test",
				Size: 100,
				Ctx:  context.Background(),
			},
			DownloadParam: struct {
				name string
				ctx  context.Context
			}{
				name: "test",
				ctx:  context.Background(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := tt.Upload.Upload(tt.UploadParam)
			if err != nil {
				t.Fatalf("Error %v\n", err)
			}
			err = tt.Download.Download(tt.DownloadParam.ctx, tt.DownloadParam.name)
			if err != nil {
				t.Fatalf("Error %v\n", err)
			}
		})
	}

}
