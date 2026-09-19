package s3Repo

import (
	"Kaban/internal/DomainLevel"
	"context"
	"errors"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Deleter struct {
	fileName string
	control  *s3.DeleteObjectInput
	ctx      context.Context
}

func NewS3Deleter() *S3Deleter {
	return &S3Deleter{}
}

func (s *S3Deleter) Deleter() error {
	_, err := s3Cred.getConnect().DeleteObject(s.ctx, s.control)
	if err != nil {
		slog.Debug("S3Deleter: error to delete an object", "ERROR", err)
		return errors.New(ErrorFileDeleted)
	}
	return nil
}

func (s *S3Deleter) SetName(s2 string) DomainLevel.MakerDeleter {
	s.fileName = s2
	return s
}

func (s *S3Deleter) SetSize(i int64) DomainLevel.MakerDeleter {
	return s
}

func (s *S3Deleter) GetFileSize() int64 {
	return 0
}

func (s *S3Deleter) GetFileName() string {

	return s.fileName
}

func (s *S3Deleter) Make(ctx context.Context) (DomainLevel.Delete, error) {
	x := &s3.DeleteObjectInput{
		Bucket: aws.String(s3Cred.getBucket()),
		Key:    aws.String(s.GetFileName()),
	}
	s.control = x
	s.ctx = ctx

	return s, nil
}
