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
	"io"
	"testing"
	"time"
)

func TestUploadEncrypt_encryptFile(t *testing.T) {
	type fields struct {
		NewUploadEncryptDataManage NewUploadEncryptDataManage
		NewUploadEncryptCrypto     NewUploadEncryptCrypto
		NewUploadEncryptDelivery   NewUploadEncryptDelivery
	}
	type args struct {
		src  io.Reader
		dst  DomainLevel.Upload
		data []byte
		ctx  context.Context
	}

	defaultFields := fields{
		NewUploadEncryptDataManage: NewUploadEncryptDataManage{
			Encode: Parsers.GetNewParsing(),
		},
		NewUploadEncryptCrypto: NewUploadEncryptCrypto{
			Generate: Crypto.GetNewGenerating(),
			Aes:      Crypto.NewAesGcmEncryption(),
			Rsa:      Crypto.NewRsaEncryption(),
		},
		NewUploadEncryptDelivery: func() NewUploadEncryptDelivery {
			flow := testFileTransfering.GetTestFlow()
			d := NewUploadEncryptDelivery{
				Uploader: flow.U,
				Deleter:  flow.D_,
			}
			return d
		}(),
	}
	defaultArgs := args{
		src: bytes.NewReader([]byte("Example data!")),
		dst: func() DomainLevel.Upload {
			u, _ := defaultFields.NewUploadEncryptDelivery.Uploader.Make(nil)
			return u
		}(),
		data: func() []byte {
			data := make([]byte, defaultFields.NewUploadEncryptCrypto.Aes.GetRequiredOverheadSize()+32)
			io.ReadFull(rand.Reader, data)
			return data
		}(),
		ctx: context.Background(),
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  defaultFields,
			args:    defaultArgs,
			wantErr: false,
		},
		{
			name:   "test2",
			fields: defaultFields,
			args: args{
				src: defaultArgs.src,
				dst: defaultArgs.dst,
				data: func() []byte {
					data := make([]byte, 32+16)
					io.ReadFull(rand.Reader, data)
					return data
				}(),
				ctx: nil,
			},
			wantErr: true,
		},
		{
			name: "test3",
			fields: fields{
				NewUploadEncryptDataManage: defaultFields.NewUploadEncryptDataManage,
				NewUploadEncryptCrypto: NewUploadEncryptCrypto{
					Aes: Crypto.NewAesCtr(),
				},
				NewUploadEncryptDelivery: defaultFields.NewUploadEncryptDelivery,
			},
			args: args{
				src: defaultArgs.src,
				dst: defaultArgs.dst,
				data: func() []byte {
					data := make([]byte, 32+12)
					io.ReadFull(rand.Reader, data)
					return data
				}(),
				ctx: nil,
			},
			wantErr: true,
		},
		{
			name:   "test4",
			fields: defaultFields,
			args: args{
				src:  defaultArgs.src,
				dst:  defaultArgs.dst,
				data: nil,
				ctx:  context.Background(),
			},
			wantErr: true,
		},
		{
			name: "test5",
			fields: fields{
				NewUploadEncryptDataManage: defaultFields.NewUploadEncryptDataManage,
				NewUploadEncryptCrypto: NewUploadEncryptCrypto{
					Generate: Crypto.GetNewGenerating(),
					Aes:      Crypto.NewAesCtr(),
					Rsa:      defaultFields.NewUploadEncryptCrypto.Rsa,
				},
				NewUploadEncryptDelivery: defaultFields.NewUploadEncryptDelivery,
			},
			args: args{
				src:  defaultArgs.src,
				dst:  defaultArgs.dst,
				data: nil,
				ctx:  defaultArgs.ctx,
			},
			wantErr: true,
		},
		{
			name: "test6",
			fields: fields{
				NewUploadEncryptDataManage: defaultFields.NewUploadEncryptDataManage,
				NewUploadEncryptCrypto: NewUploadEncryptCrypto{
					Generate: defaultFields.NewUploadEncryptCrypto.Generate,
					Aes:      Crypto.NewAesGcmError(),
				},
				NewUploadEncryptDelivery: defaultFields.NewUploadEncryptDelivery,
			},
			args: args{
				src: defaultArgs.src,
				dst: defaultArgs.dst,
				data: func() []byte {
					key := make([]byte, 32+defaultFields.NewUploadEncryptCrypto.Aes.GetRequiredOverheadSize())
					io.ReadFull(rand.Reader, key)
					return key
				}(),
				ctx: defaultArgs.ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sa := &UploadEncrypt{
				NewUploadEncryptDataManage: tt.fields.NewUploadEncryptDataManage,
				NewUploadEncryptCrypto:     tt.fields.NewUploadEncryptCrypto,
				NewUploadEncryptDelivery:   tt.fields.NewUploadEncryptDelivery,
			}
			if err := sa.encryptFile(tt.args.src, tt.args.dst, tt.args.data, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("encryptFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestUploadEncrypt_UploadEncryptApplication(t *testing.T) {
	type fields struct {
		NewUploadEncryptDataManage NewUploadEncryptDataManage
		NewUploadEncryptCrypto     NewUploadEncryptCrypto
		Keys                       DomainLevel.CryptoKeys
		NewUploadEncryptDelivery   NewUploadEncryptDelivery
	}
	type args struct {
		makeArgs func() args
		data     UploadEncryptIncomeData
	}

	defaultFields := fields{
		NewUploadEncryptDataManage: NewUploadEncryptDataManage{
			Encode: Parsers.GetNewParsing(),
		},
		NewUploadEncryptCrypto: NewUploadEncryptCrypto{
			Generate: Crypto.GetNewGenerating(),
			Aes:      Crypto.NewAesGcmEncryption(),
			Rsa:      Crypto.NewRsaEncryption(),
		},
		Keys: func() DomainLevel.CryptoKeys {
			DomainLevel.GetCryptoKeys = RepoEncrypterKeys.NewNormalKeyTest()
			return DomainLevel.GetCryptoKeys
		}(),
		NewUploadEncryptDelivery: func() NewUploadEncryptDelivery {
			flow := testFileTransfering.GetTestFlow()
			d := NewUploadEncryptDelivery{
				Uploader:     flow.U,
				RedisWriter:  RedisInteration.NewRedisWriteTest(),
				RedisChecker: RedisInteration.NewRedisCheckTest(),
				RedisDeleter: RedisInteration.NewRedisDeleteTest(),
				Deleter:      flow.D_,
			}
			return d
		}(),
	}

	defaultArgs := args{
		makeArgs: func() args {
			d := UploadEncryptIncomeData{
				File: bytes.NewReader([]byte("Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!Hello!!")),
				Name: []byte("Hello"),
				Ctx:  context.Background(),
				Size: 100,
			}
			return args{
				data: d,
			}
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  defaultFields,
			args:    defaultArgs.makeArgs(),
			wantErr: false,
		},
		{
			name:   "test2",
			fields: defaultFields,
			args: args{
				data: UploadEncryptIncomeData{
					File: defaultArgs.data.File,
					Name: defaultArgs.data.Name,
					Ctx:  defaultArgs.data.Ctx,
					Size: 1000 * 1024 * 1024,
				},
			},
			wantErr: true,
		},
		{
			name: "test3",
			fields: fields{
				NewUploadEncryptDataManage: defaultFields.NewUploadEncryptDataManage,
				NewUploadEncryptCrypto: NewUploadEncryptCrypto{
					Generate: defaultFields.NewUploadEncryptCrypto.Generate,
					Aes:      Crypto.NewAesGcmError(),
					Rsa:      defaultFields.NewUploadEncryptCrypto.Rsa,
				},
				Keys:                     defaultFields.Keys,
				NewUploadEncryptDelivery: defaultFields.NewUploadEncryptDelivery,
			},
			args:    defaultArgs.makeArgs(),
			wantErr: true,
		},
		{
			name: "test4",
			fields: fields{
				NewUploadEncryptDataManage: defaultFields.NewUploadEncryptDataManage,
				NewUploadEncryptCrypto: NewUploadEncryptCrypto{
					Generate: defaultFields.NewUploadEncryptCrypto.Generate,
					Aes:      defaultFields.NewUploadEncryptCrypto.Aes,
					Rsa:      Crypto.NewRsaError(),
				},
				Keys:                     defaultFields.Keys,
				NewUploadEncryptDelivery: defaultFields.NewUploadEncryptDelivery,
			},
			args:    defaultArgs.makeArgs(),
			wantErr: true,
		},
		{
			name: "test5",
			fields: fields{
				NewUploadEncryptDataManage: defaultFields.NewUploadEncryptDataManage,
				NewUploadEncryptCrypto:     defaultFields.NewUploadEncryptCrypto,
				Keys:                       defaultFields.Keys,
				NewUploadEncryptDelivery: NewUploadEncryptDelivery{
					Uploader:     defaultFields.NewUploadEncryptDelivery.Uploader,
					RedisWriter:  RedisInteration.NewRedisWriteErrorTest(),
					RedisChecker: defaultFields.NewUploadEncryptDelivery.RedisChecker,
					RedisDeleter: defaultFields.NewUploadEncryptDelivery.RedisDeleter,
					Deleter:      defaultFields.NewUploadEncryptDelivery.Deleter,
				},
			},
			args:    defaultArgs.makeArgs(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UploadEncrypt{
				NewUploadEncryptDataManage: tt.fields.NewUploadEncryptDataManage,
				NewUploadEncryptCrypto:     tt.fields.NewUploadEncryptCrypto,
				NewUploadEncryptDelivery:   tt.fields.NewUploadEncryptDelivery,
			}
			_, err := u.UploadEncryptApplication(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadEncryptApplication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			time.Sleep(10 * time.Millisecond)
		})
	}
}
