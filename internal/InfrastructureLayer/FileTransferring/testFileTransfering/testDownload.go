package testFileTransfering

import (
	"Kaban/internal/DomainLevel"
	"context"
	"io"
	"os"
)

type TestDownload struct {
	f *os.File
}

func NewTestDownload() *TestDownload {
	return &TestDownload{}
}

func (t TestDownload) Downloader() (io.Reader, error) {
	return t.f, nil
}

func (t TestDownload) CloseSource() error {
	return t.f.Close()
}

func (t TestDownload) SetName(s string) DomainLevel.MakerDownloader {
	//TODO implement me
	panic("implement me")
}

func (t TestDownload) SetSize(i int64) DomainLevel.MakerDownloader {
	//TODO implement me
	panic("implement me")
}

func (t TestDownload) GetFileSize() int64 {
	//TODO implement me
	panic("implement me")
}

func (t TestDownload) GetFileName() string {
	//TODO implement me
	panic("implement me")
}

func (t *TestDownload) Make(ctx context.Context) (DomainLevel.Download, error) {
	file, err := os.Open("example1.txt")
	if err != nil {
		panic(err)
	}

	t.f = file

	return t, nil
}

func (t TestDownload) SetAdditionalReader(reader io.Reader) DomainLevel.MakerDownloader {
	//TODO implement me
	panic("implement me")
}
