package s3Repo

import (
	"Kaban/internal/DomainLevel"
	"context"
	"errors"

	"log/slog"

	NewVersion "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Download struct {
	S3Info Variables
}

func GetNewS3Download(s3Info Variables) *S3Download {
	return &S3Download{S3Info: s3Info}
}

type DownloadingS3 interface {
	GetDownload(string, context.Context) (*NewVersion.GetObjectOutput, error)
	GetDownloadSecure(context.Context, string) (*s3.GetObjectOutput, error)
}

func (s S3Download) GetDownloadSecure(ctx context.Context, name string) (*s3.GetObjectOutput, error) {

	dsa := s3.New(s.S3Info.OldConnect)
	Params := &s3.GetObjectInput{Bucket: aws.String(s.S3Info.Bucket), Key: aws.String(name)}

	O, err := dsa.GetObjectWithContext(ctx, Params)
	if err != nil {
		slog.Error("GetDownloadSecure; the file wasn't find", "ERROR", err)
		return nil, errors.New(DomainLevel.ErrorCantFindFile)
	}

	return O, nil
}

func (s S3Download) GetDownload(TrueFileName string, ctx context.Context) (*NewVersion.GetObjectOutput, error) {

	InputData := &NewVersion.GetObjectInput{Bucket: aws.String(s.S3Info.Bucket), Key: aws.String(TrueFileName)}

	S, err := s.S3Info.S3Connect.GetObject(ctx, InputData, func(options *NewVersion.Options) {
		options.DisableLogOutputChecksumValidationSkipped = true

	})

	if err != nil {
		slog.Error("S3 GetDownload; the error happened", "ERROR", err.Error())
		return nil, errors.New(DomainLevel.ErrorStartDownloading)
	}

	return S, err

}
