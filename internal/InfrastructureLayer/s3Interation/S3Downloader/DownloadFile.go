package S3Downloader

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func (s S3Download) Download(TrueFileName string, ctx context.Context) (*s3.GetObjectOutput, error) {

	InputData := &s3.GetObjectInput{Bucket: aws.String(s.S3Info.Bucket), Key: aws.String(TrueFileName)}

	S, err := s.S3Info.S3Connect.GetObject(ctx, InputData, func(options *s3.Options) {
		options.DisableLogOutputChecksumValidationSkipped = true

	})

	if err != nil {
		slog.Error("Error getting s3 object", "Error", err.Error())
		return nil, err
	}

	return S, err

}
