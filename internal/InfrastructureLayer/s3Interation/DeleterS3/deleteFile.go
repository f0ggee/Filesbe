package DeleterS3

import (
	"context"
	"log/slog"

	"Kaban/internal/InfrastructureLayer/s3Interation"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func (d *DeleterS3) DeleteFileFromS3(key string, ctx context.Context) error {
	s := &s3.DeleteObjectInput{
		Bucket: aws.String(s3Interation.S3Info.Bucket),
		Key:    &key,
	}

	slog.Info("Bucket name", s3Interation.S3Info.Bucket)
	_, err := d.Conf.DeleteObject(ctx, s)
	if err != nil {
		slog.Error("Error in delete func", err)
		return err
	}
	return nil
}
