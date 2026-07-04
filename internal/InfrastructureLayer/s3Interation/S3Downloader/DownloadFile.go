package S3Downloader

import (
	"Kaban/internal/DomainLevel"
	"context"
	"errors"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func (s S3Download) GetDownload(TrueFileName string, ctx context.Context) (*s3.GetObjectOutput, error) {

	InputData := &s3.GetObjectInput{Bucket: aws.String(s.S3Info.Bucket), Key: aws.String(TrueFileName)}

	S, err := s.S3Info.S3Connect.GetObject(ctx, InputData, func(options *s3.Options) {
		options.DisableLogOutputChecksumValidationSkipped = true

	})

	if err != nil {
		slog.Error("S3 GetDownload; the error happened", "ERROR", err.Error())
		return nil, errors.New(DomainLevel.ErrorStartDownloading)
	}

	return S, err

}
