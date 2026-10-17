package testFileTransfering

import (
	"Kaban/internal/DomainLevel"
	"bytes"
	"context"
	"io"
)

type TestDownloader struct {
}

func GetTestDownloader() *TestDownloader {
	return &TestDownloader{}
}

func (t TestDownloader) Downloader() (io.Reader, error) {
	return bytes.NewReader(Hub), nil
}

func (t TestDownloader) CloseSource() error {
	return nil
}

func (t TestDownloader) SetName(s string) DomainLevel.MakerDownloader {
	return t
}

func (t TestDownloader) SetSize(i int64) DomainLevel.MakerDownloader {
	return t
}

func (t TestDownloader) GetFileSize() int64 {
	return 100
}

func (t TestDownloader) GetFileName() string {
	return "test"
}

func (t TestDownloader) Make(ctx context.Context) (DomainLevel.Download, error) {
	return t, nil
}

func (t TestDownloader) SetAdditionalReader(reader io.Reader) DomainLevel.MakerDownloader {
	return t
}
