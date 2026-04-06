package S3Downloader

import (
	"Kaban/internal/InfrastructureLayer/s3Interation"
	"context"
	"io"
)

type S3Downloading struct {
	S3Info s3Interation.Variables
}

func (s3 S3Downloading) DownloadSecure(ctx context.Context, name string, writer *io.PipeWriter, aesKey []byte, realFileName string, Reader *io.PipeReader) (*io.PipeReader, error) {
	//TODO implement me
	panic("implement me")
}
