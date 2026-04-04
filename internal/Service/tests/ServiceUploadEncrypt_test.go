package Handlers

import (
	"Kaban/internal/InfrastructureLayer/RedisInteration/DeletingRedis"
	"Kaban/internal/InfrastructureLayer/RedisInteration/RedisChecking"
	"Kaban/internal/InfrastructureLayer/s3Interation/DeleterS3"
	Handlers2 "Kaban/internal/Service/Handlers"
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const RedisChecker = "IsFllen"
const DeleterFileInfo = "isFallRedis"
const S3DeleterFile = "IsFall"

func TestTester(t *testing.T) {
	type fields struct {
		S3                  Handlers2.S3Controlling
		Crypto              Handlers2.HandlerPackCrypto
		FileInfo            Handlers2.HandlerFileManagerPack
		AuthTokens          Handlers2.HandlerPackAuthTokens
		DatabaseControlling Handlers2.DatabaseControlling
		RedisControlling    Handlers2.RedisControlling
		Grpc                Handlers2.HandlerGrpc
		Convert             Handlers2.Converter
		Keys                Handlers2.KeysControlling
	}
	type args struct {
		shortNameForFile string
		ctx              context.Context
	}

	s3 := Handlers2.S3Controlling{
		Deleter: &DeleterS3.DeleterS3{Conf: nil},
	}
	Redis := Handlers2.RedisControlling{
		CheckerRedis: &RedisChecking.ValidationRedis{Re: nil},
		Deleter:      &DeletingRedis.DeleterRedis{Re: nil},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Tester1",

			args: args{

				ctx:              context.WithValue(context.Background(), RedisChecker, true),
				shortNameForFile: "Tester1",
			},
			wantErr: true,
			fields: fields{
				S3:               s3,
				RedisControlling: Redis,
			},
		},
		{
			name: "Tester2",
			args: args{
				ctx:              context.WithValue(context.Background(), S3DeleterFile, true),
				shortNameForFile: "Tester2",
			},
			wantErr: true,
			fields: fields{
				S3:               s3,
				RedisControlling: Redis,
			},
		},
		{
			name: "Tester3",
			args: args{
				ctx:              context.Background(),
				shortNameForFile: "Tester3",
			},
			wantErr: false,
			fields: fields{
				S3:               s3,
				RedisControlling: Redis,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sa := &Handlers2.HandlerPackCollect{
				S3:                  tt.fields.S3,
				Crypto:              tt.fields.Crypto,
				FileInfo:            tt.fields.FileInfo,
				AuthTokens:          tt.fields.AuthTokens,
				DatabaseControlling: tt.fields.DatabaseControlling,
				RedisControlling:    tt.fields.RedisControlling,
				Grpc:                tt.fields.Grpc,
				Convert:             tt.fields.Convert,
				Keys:                tt.fields.Keys,
			}
			err := sa.Tester(tt.args.shortNameForFile, tt.args.ctx)
			if (err == nil) && tt.wantErr || err != nil && !tt.wantErr {
				t.Errorf("We expected error wanted %v from %v but got  %v", tt.wantErr, tt.name, err)
			} else {
				t.Logf("%v %v %v "+
					"passed", tt.name, tt.wantErr, err)
			}
		})
	}
}

func Test_uploadFileEncrypt(t *testing.T) {
	type args struct {
		cfg           *s3.Client
		BesParts      int
		goroutine     int
		ctx           context.Context
		shortFileName string
		ContentType   string
		reader        *io.PipeReader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := uploadFileEncrypt(tt.args.cfg, tt.args.BesParts, tt.args.goroutine, tt.args.ctx, tt.args.shortFileName, tt.args.ContentType, tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("uploadFileEncrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("uploadFileEncrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
