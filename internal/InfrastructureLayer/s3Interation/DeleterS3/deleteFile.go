package DeleterS3

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func (d *DeleterS3) DeleteFileFromS3(key string, ctx context.Context) error {
	s := &s3.DeleteObjectInput{
		Bucket: aws.String(d.S3Info.Bucket),
		Key:    &key,
	}
	_, err := d.S3Info.S3Connect.DeleteObject(ctx, s)
	if err != nil {
		slog.Error("Error in delete func", "ERROR", err.Error())
		return err
	}
	return nil
}
