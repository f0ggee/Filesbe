package s3Repo

import (
	"Kaban/internal/DomainLevel"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type NewDeleterS3 struct {
	S3Info Variables
}

type DeleterS3 interface {
	DeleteFileFromS3(string, context.Context) error
	DeleterS3Test(string, context.Context) error
}

func GetNewDeleterS3(s3Info Variables) *NewDeleterS3 {
	return &NewDeleterS3{S3Info: s3Info}
}

func (d *NewDeleterS3) DeleterS3Test(s string, Cont context.Context) error {

	time.Sleep(2 * time.Second)
	sa, de := Cont.Value("IsFall").(bool)
	if sa != false {
		if de {
			return errors.New("error by s3")
		}
	}
	return nil

}
func (d *NewDeleterS3) DeleteFileFromS3(key string, ctx context.Context) error {
	s := &s3.DeleteObjectInput{
		Bucket: aws.String(d.S3Info.Bucket),
		Key:    &key,
	}
	_, err := d.S3Info.S3Connect.DeleteObject(ctx, s)
	if err != nil {
		slog.Error("DeleteFileFromS3; error to delete a file", "ERROR", err)
		return errors.New(DomainLevel.ErrorFilNotDeleted)
	}
	return nil
}
