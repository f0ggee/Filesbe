package Application

import (
	"Kaban/internal/InfrastructureLayer/Crypto"
	"Kaban/internal/InfrastructureLayer/FileTransferring/testFileTransfering"
	"Kaban/internal/InfrastructureLayer/Parsers"
	"Kaban/internal/InfrastructureLayer/RedisInteration"
	"bytes"
	"context"
	"net/http/httptest"
	"testing"
)

func TestApplicationFile_appFile(t *testing.T) {
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
			name: "test2",
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
			_, err := tt.Upload.UploadApplication(tt.UploadParam)
			if err != nil {
				t.Fatalf("Error %v\n", err)
			}
			err = tt.Download.DownloadApplication(tt.DownloadParam.ctx, tt.DownloadParam.name)
			if err != nil {
				t.Fatalf("Error %v\n", err)
			}
		})
	}

}
