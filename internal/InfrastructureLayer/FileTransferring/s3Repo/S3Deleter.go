package s3Repo

import (
	"Kaban/internal/DomainLevel"
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Delete struct {
	fileName string
	control  *s3.DeleteObjectInput
	ctx      context.Context
}

func NewS3Deleter() S3Delete {
	return S3Delete{}
}

func (s *S3Delete) Deleter() error {
	_, err := s3Cred.getConnect().DeleteObject(s.ctx, s.control)
	if err != nil {
		slog.Debug("S3Delete: error to delete an object", "ERROR", err)
		return ErrorFileDeleted
	}
	return nil
}

func (s *S3Delete) SetName(s2 string) DomainLevel.MakerDeleter {
	s.fileName = s2
	return s
}

func (s *S3Delete) SetSize(i int64) DomainLevel.MakerDeleter {
	return s
}

func (s *S3Delete) GetFileSize() int64 {
	return 0
}

func (s *S3Delete) GetFileName() string {

	return s.fileName
}

func (s *S3Delete) Make(ctx context.Context) (DomainLevel.Delete, error) {
	x := &s3.DeleteObjectInput{
		Bucket: aws.String(s3Cred.getBucket()),
		Key:    aws.String(s.GetFileName()),
	}
	s.control = x
	s.ctx = ctx

	return s, nil
}
