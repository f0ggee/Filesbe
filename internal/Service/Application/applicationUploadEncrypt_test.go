package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Drivers/Crypto"
	"Kaban/internal/Drivers/FileTransferring/testFileTransfering"
	"Kaban/internal/Drivers/Parsers"
	"Kaban/internal/Drivers/RedisInteration"
	"Kaban/internal/Drivers/RepoEncrypterKeys"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"testing"
)

func TestUploadEncrypt_UploadEncryptApplication(t *testing.T) {
	rsaK, _ := rsa.GenerateKey(rand.Reader, 2048)

	DomainLevel.GetCryptoKeys = RepoEncrypterKeys.NewTestEncrypterKeys(x509.MarshalPKCS1PrivateKey(rsaK))
	type fields struct {
		NewUploadEncryptDataManage NewUploadEncryptDataManage
		NewUploadEncryptCrypto     NewUploadEncryptCrypto
		NewUploadEncryptDelivery   NewUploadEncryptDelivery
	}
	newUpload := NewUploadEncryptDelivery{
		Uploader:     nil,
		RedisWriter:  RedisInteration.NewRedisWriteTest(),
		RedisChecker: RedisInteration.NewRedisCheckTest(),
		RedisDeleter: RedisInteration.NewRedisDeleteTest(),
		Deleter:      nil,
	}
	crypto := NewUploadEncryptCrypto{
		Generate: Crypto.GetNewGenerating(),
		Aes:      Crypto.NewAesCtr(),
		Rsa:      Crypto.NewRsaEncryption(),
	}
	parser := NewUploadEncryptDataManage{Parsers.GetNewParsing()}
	type args struct {
		data UploadEncryptIncomeData
	}
	fileContent := []byte("Hello!")
	testArgs := args{
		data: UploadEncryptIncomeData{
			File: bytes.NewReader(fileContent),
			Name: []byte("test"),
			Ctx:  context.Background(),
			Size: 100,
		},
	}
	bigFileArgs := testArgs
	bigFileArgs.data.Size = 600 * 1024 * 1024
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		flow    testFileTransfering.Flow
	}{
		{
			flow: *testFileTransfering.GetTestFlow(),
			name: "test1",
			fields: fields{
				NewUploadEncryptDataManage: NewUploadEncryptDataManage{
					Encode: Parsers.GetNewParsing(),
				},
				NewUploadEncryptCrypto:   crypto,
				NewUploadEncryptDelivery: newUpload,
			},
			args:    testArgs,
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				NewUploadEncryptDataManage: parser,
				NewUploadEncryptCrypto:     crypto,
				NewUploadEncryptDelivery:   newUpload,
			},
			args:    bigFileArgs,
			wantErr: true,
		},
		{
			name: "test3",
			fields: fields{
				NewUploadEncryptDataManage: parser,
				NewUploadEncryptCrypto:     crypto,
				NewUploadEncryptDelivery:   newUpload,
			},
			args:    testArgs,
			wantErr: true,
			flow:    testFileTransfering.GetTestErrorFlow(),
		},
		{
			name: "test4",
			fields: fields{
				NewUploadEncryptDataManage: parser,
				NewUploadEncryptCrypto: NewUploadEncryptCrypto{
					Generate: crypto.Generate,
					Aes:      crypto.Aes,
					Rsa:      Crypto.NewRsaError(),
				},
				NewUploadEncryptDelivery: newUpload,
			},
			args:    testArgs,
			wantErr: true,
			flow:    *testFileTransfering.GetTestFlow(),
		},
		{
			name: "test5",
			fields: fields{
				NewUploadEncryptDataManage: parser,
				NewUploadEncryptCrypto: NewUploadEncryptCrypto{
					Generate: crypto.Generate,
					Aes:      Crypto.NewAesCtrError{},
					Rsa:      crypto.Rsa,
				},
				NewUploadEncryptDelivery: newUpload,
			},
			args:    testArgs,
			wantErr: true,
			flow:    *testFileTransfering.GetTestFlow(),
		},
		{
			name: "test6",
			fields: fields{
				NewUploadEncryptDataManage: parser,
				NewUploadEncryptCrypto:     crypto,
				NewUploadEncryptDelivery:   newUpload,
			},
			args:    testArgs,
			wantErr: true,
			flow:    testFileTransfering.GetTestErrorFlow(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// every case shares testArgs, so give each run its own unread reader
			tt.args.data.File = bytes.NewReader(fileContent)
			tt.fields.NewUploadEncryptDelivery.Uploader = tt.flow.U
			tt.fields.NewUploadEncryptDelivery.Deleter = tt.flow.D_
			sa := &UploadEncrypt{
				NewUploadEncryptDataManage: tt.fields.NewUploadEncryptDataManage,
				NewUploadEncryptCrypto:     tt.fields.NewUploadEncryptCrypto,
				NewUploadEncryptDelivery:   tt.fields.NewUploadEncryptDelivery,
			}
			_, err := sa.UploadEncryptApplication(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadEncryptApplication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
