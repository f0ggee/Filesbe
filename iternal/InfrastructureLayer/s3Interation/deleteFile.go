package s3Interation

import (
	"Kaban/iternal/Service/Helpers"
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type ConntrolerForS3 struct{}

func (*ConntrolerForS3) DeleteFileFromS3(key string, bucket string) error {
	ConfigureS3, err := Helpers.S3Helper()
	if err != nil {
		slog.Error("can't connect to S3 server", "Err", err)
		return err
	}
	s := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    &key,
	}

	_, err = ConfigureS3.DeleteObject(context.Background(), s)
	if err != nil {
		slog.Error("Error in delete func", err)
		return err
	}
	return nil
}
