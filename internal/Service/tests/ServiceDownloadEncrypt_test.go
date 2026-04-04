package Handlers

import (
	Handlers2 "Kaban/internal/Service/Handlers"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestDecryptFile(t *testing.T) {
	type args struct {
		AesKey []byte
		o      *s3.GetObjectOutput
		writer *io.PipeWriter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Handlers2.DecryptFile(tt.args.AesKey, tt.args.o, tt.args.writer); (err != nil) != tt.wantErr {
				t.Errorf("DecryptFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandlerPackCollect_DownloadEncrypt(t *testing.T) {
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
		w    http.ResponseWriter
		ctxs context.Context
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := sa.DownloadEncrypt(tt.args.w, tt.args.ctxs, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DownloadEncrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandlerPackCollect_downloadFileToClient(t *testing.T) {
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
		w            http.ResponseWriter
		ctx          context.Context
		name         string
		writer       *io.PipeWriter
		aesKey       []byte
		realFileName string
		Reader       *io.PipeReader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := sa.downloadFileToClient(tt.args.w, tt.args.ctx, tt.args.name, tt.args.writer, tt.args.aesKey, tt.args.realFileName, tt.args.Reader); (err != nil) != tt.wantErr {
				t.Errorf("downloadFileToClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
