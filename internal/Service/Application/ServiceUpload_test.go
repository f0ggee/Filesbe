package Application

import (
	"Kaban/internal/InfrastructureLayer/Crypto"
	"Kaban/internal/InfrastructureLayer/FileTransferring/testFileTransfering"
	"Kaban/internal/InfrastructureLayer/Parsers"
	"Kaban/internal/InfrastructureLayer/RedisInteration"
	"bytes"
	"context"
	"testing"
)

func TestNewUpload_UploadApplication(t *testing.T) {
	type fields struct {
		NewFileUploaderCrypto    NewFileUploaderCrypto
		NewFileUploaderDataMange NewFileUploaderDataMange
		NewFileUploaderDelivery  NewFileUploaderDelivery
	}
	type args struct {
		r UploadIncomeData
	}

	inputData := bytes.NewReader([]byte("Hello!"))
	defaultFields := fields{
		NewFileUploaderCrypto: NewFileUploaderCrypto{
			Generator: Crypto.GetNewGenerating(),
		},
		NewFileUploaderDataMange: NewFileUploaderDataMange{
			Encode: Parsers.GetNewParsing(),
		},
		NewFileUploaderDelivery: NewFileUploaderDelivery{
			Uploader:   testFileTransfering.GetTestFlow().U,
			WriteRedis: RedisInteration.NewRedisWriteTest(),
		},
	}
	defaultParams := args{
		UploadIncomeData{
			File: inputData,
			Name: "12345",
			Size: 100,
			Ctx:  context.Background(),
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
			args:    defaultParams,
			wantErr: false,
		},
		{
			name:   "test2",
			fields: defaultFields,
			args: args{
				r: UploadIncomeData{
					File: defaultParams.r.File,
					Name: defaultParams.r.Name,
					Size: 700 * 1024 * 1024,
					Ctx:  defaultParams.r.Ctx,
				},
			},
			wantErr: true,
		},
		{
			name: "test3",
			fields: fields{
				NewFileUploaderCrypto:    defaultFields.NewFileUploaderCrypto,
				NewFileUploaderDataMange: defaultFields.NewFileUploaderDataMange,
				NewFileUploaderDelivery: NewFileUploaderDelivery{
					Uploader:   testFileTransfering.GetTestErrorFlow().U,
					WriteRedis: RedisInteration.NewRedisWriteTest(),
				},
			},
			args:    defaultParams,
			wantErr: true,
		},
		{
			name: "test4",
			fields: fields{
				NewFileUploaderCrypto: defaultFields.NewFileUploaderCrypto,
				NewFileUploaderDataMange: NewFileUploaderDataMange{
					Encode: Parsers.GetParseError(),
				},
				NewFileUploaderDelivery: defaultFields.NewFileUploaderDelivery,
			},
			args:    defaultParams,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sa := &NewUpload{
				NewFileUploaderCrypto:    tt.fields.NewFileUploaderCrypto,
				NewFileUploaderDataMange: tt.fields.NewFileUploaderDataMange,
				NewFileUploaderDelivery:  tt.fields.NewFileUploaderDelivery,
			}
			_, err := sa.UploadApplication(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadApplication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
