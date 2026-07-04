package S3Downloader

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/s3Interation"
	"context"
	"errors"

	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Download struct {
	S3Info s3Interation.Variables
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
