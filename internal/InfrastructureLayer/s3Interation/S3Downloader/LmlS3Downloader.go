package S3Downloader

import (
	"Kaban/internal/InfrastructureLayer/s3Interation"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Download struct {
	S3Info s3Interation.Variables
}

func (s S3Download) DownloadSecure(ctx context.Context, name string) (io.ReadCloser, int64, error) {

	dsa := s3.New(s.S3Info.OldConnect)
	Params := &s3.GetObjectInput{Bucket: aws.String(s.S3Info.Bucket), Key: aws.String(name)}

	O, err := dsa.GetObjectWithContext(ctx, Params)
	if err != nil {
		slog.Error("Error finding the file in the s3 storage", "Err", err)
		return nil, 0, err
	}

	switch {
	case strings.Contains(fmt.Sprint(err), "NoSuchKey"):
		slog.Info("File was used")
		return nil, 0, errors.New("file was used")

	case errors.Is(err, context.DeadlineExceeded):
		slog.Error("Time was exceeded")
		return nil, 0, errors.New("time was exceeded")
	case errors.Is(err, context.Canceled):
		slog.Info("a user has been cancelled download ")
		return nil, 0, errors.New("a user has been canceled download ")

	}
	if err != nil {
		slog.Error("ServiceDownload:", "Error", err.Error())
		return nil, 0, err
	}

	return O.Body, *O.ContentLength, nil
}
