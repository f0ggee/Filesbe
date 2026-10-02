package s3Repo

import (
	"Kaban/internal/DomainLevel"
	"context"
	"crypto/rand"
	"errors"
	"io"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
)

type S3Upload struct {
	uploader  *transfermanager.Client
	ctx       context.Context
	extension string
	name      string
	size      int64
}

func (s *S3Upload) SetName(s2 string) DomainLevel.MakerUploader {
	s.name = s2
	return s
}

func (s *S3Upload) SetSize(i int64) DomainLevel.MakerUploader {

	s.size = i
	return s
}

func (s *S3Upload) SetAdditionalWriter(writer io.Writer) DomainLevel.MakerUploader {
	return nil
}
func NewS3Upload() *S3Upload {
	return &S3Upload{}
}
func (s *S3Upload) Uploader(reader io.Reader) error {

	sa, err := s.uploader.UploadObject(s.ctx, &transfermanager.UploadObjectInput{
		Bucket:      aws.String(s3Cred.getBucket()),
		Key:         aws.String(s.name),
		Body:        reader,
		ContentType: aws.String(s.extension),
	})
	if err == nil {
		s.size = *sa.ContentLength
		return nil
	}

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		slog.Error("S3 uploader: a client stopped uploading")
		return ErrorClientStop

	case err != nil:
		slog.Error("S3 uploader: a strange error", "ERROR", err)
		return ErrorStrangeError
	}
	return ErrorStrangeError
}

func (s *S3Upload) CloseSource() error {
	return nil
}

func (s *S3Upload) GetFileSize() int64 {

	return s.size
}

func (s *S3Upload) GetFileName() string {
	return s.name
}

func (s *S3Upload) Make(ctx context.Context) (DomainLevel.Upload, error) {
	if s.size == 0 {
		return nil, errors.New(ErrorFileSize)
	}
	if s.name == "" {
		var x strings.Builder
		x.WriteString("NotIdentifyFile_")
		x.WriteString(rand.Text()[:10])
		s.SetName(x.String())
	}
	settings := DomainLevel.GetNewFileSettings(s.size, s.name)
	parts, goroutines := settings.FindBestOptions()
	s.extension = settings.FindFormatOfFile()
	client := transfermanager.New(s3Cred.getConnect(), func(options *transfermanager.Options) {
		options.Concurrency = goroutines
		options.PartSizeBytes = int64(parts)
	})
	s.ctx = ctx
	s.uploader = client
	return s, nil
}
