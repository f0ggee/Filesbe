package s3Repo

import (
	"Kaban/internal/DomainLevel"
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"log/slog"
)

type S3Downloader struct {
	name string
	size int64
	ctx  context.Context
	obj  *transfermanager.Client
	body *io.Reader
}

func (s *S3Downloader) SetAdditionalReader(reader io.Reader) DomainLevel.MakerDownloader {
	return s
}

func NewS3Downloader() *S3Downloader {
	return &S3Downloader{}
}

func (s *S3Downloader) Downloader() (io.Reader, error) {
	output, err := s.obj.GetObject(s.ctx, &transfermanager.GetObjectInput{
		Bucket: aws.String(s3Cred.getBucket()),
		Key:    aws.String(s.name),
	})
	if err == nil {
		s.body = &output.Body
		s.SetSize(*output.ContentLength)
		return output.Body, nil
	}
	var ns *types.NoSuchKey
	switch {
	case errors.Is(err, ns):
		slog.Error("S3Downloader: the file wasn't found")
		return nil, errors.New(ErrorNoFile)

	case err != nil:
		slog.Error("S3Downloader: an strange error happened during downloading", "ERROR", err)
		return nil, ErrorStrangeError
	}
	return nil, ErrorStrangeError
}

func (s *S3Downloader) CloseSource() error {
	return nil
}

func (s *S3Downloader) SetName(s2 string) DomainLevel.MakerDownloader {
	s.name = s2
	return s
}

func (s *S3Downloader) SetSize(i int64) DomainLevel.MakerDownloader {
	s.size = i
	return s
}

func (s *S3Downloader) GetFileSize() int64 {
	return s.size
}

func (s *S3Downloader) GetFileName() string {
	return s.name
}

func (s *S3Downloader) Make(ctx context.Context) (DomainLevel.Download, error) {
	if s.size == 0 {
		slog.Error("S3 downloader: the file size isn't defined")
		return nil, errors.New(ErrorFileSize)
	}
	if s.name == "" {
		slog.Error("S3 downloader: the file name isn't defined")
		return nil, errors.New(ErrorFileName)
	}

	opts, goroutines := DomainLevel.GetNewFileSettings(s.size, "").FindBestOptions()
	downloadObj := transfermanager.New(s3Cred.getConnect(), func(options *transfermanager.Options) {
		options.Concurrency = goroutines
		options.PartSizeBytes = int64(opts)
	})
	s.ctx = ctx

	s.obj = downloadObj
	return s, nil
}
