package Application

import (
	"Kaban/internal/InfrastructureLayer/Crypto"
	"Kaban/internal/InfrastructureLayer/Parsers"
	"Kaban/internal/InfrastructureLayer/RedisInteration"
	"testing"
)

func TestNewUploadEncrypt_UploadEncrypt(t *testing.T) {
	type fields struct {
		NewUploadEncryptDataManage NewUploadEncryptDataManage
		NewUploadEncryptCrypto     NewUploadEncryptCrypto
		NewUploadEncryptDelivery   NewUploadEncryptDelivery
	}
	type args struct {
		data IncomeData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				NewUploadEncryptDataManage: NewUploadEncryptDataManage{
					Encode: Parsers.GetNewParsing(),
				},
				NewUploadEncryptCrypto: NewUploadEncryptCrypto{
					Generate: Crypto.GetNewGenerating(),
					Encrypt:  Crypto.NewAesGcmEncryption(),
				},
				NewUploadEncryptDelivery: NewUploadEncryptDelivery{
					Uploader:     nil,
					RedisWriter:  RedisInteration.NewRedisWriteTest(),
					RedisChecker: RedisInteration.NewRedisCheckTest(),
					RedisDeleter: RedisInteration.NewNewRedisDeleteTest(),
					Deleter:      nil,
				},
			},
			args: args{
				data: IncomeData{
					File: nil,
					Name: "",
					Ctx:  nil,
					Size: 0,
				},
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sa := &NewUploadEncrypt{
				NewUploadEncryptDataManage: tt.fields.NewUploadEncryptDataManage,
				NewUploadEncryptCrypto:     tt.fields.NewUploadEncryptCrypto,
				NewUploadEncryptDelivery:   tt.fields.NewUploadEncryptDelivery,
			}
			got, err := sa.UploadEncrypt(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadEncrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UploadEncrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
