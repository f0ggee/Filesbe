package testFileTransfering

import (
	"Kaban/internal/DomainLevel"
	"bytes"
	"context"
	"io"
)

type testDownload struct {
	h *[]byte
}

func (t testDownload) Downloader() (io.Reader, error) {
	return bytes.NewReader(*t.h), nil
}

func (t testDownload) CloseSource() error {
	return nil
}

func (t testDownload) SetName(s string) DomainLevel.MakerDownloader {
	return t
}

func (t testDownload) SetSize(i int64) DomainLevel.MakerDownloader {
	return t
}

func (t testDownload) GetFileSize() int64 {
	return 100
}

func (t testDownload) GetFileName() string {
	return "data"
}

func (t testDownload) Make(ctx context.Context) (DomainLevel.Download, error) {
	return t, nil
}

func (t testDownload) SetAdditionalReader(reader io.Reader) DomainLevel.MakerDownloader {
	return t
}

type DownloadError struct {
	err error
}

func NewDownloadError(err error) DownloadError {
	return DownloadError{err: err}
}

func (d DownloadError) Downloader() (io.Reader, error) {
	return nil, d.err
}

func (d DownloadError) CloseSource() error {
	return nil
}

func (d DownloadError) SetName(s string) DomainLevel.MakerDownloader {
	return d
}

func (d DownloadError) SetSize(i int64) DomainLevel.MakerDownloader {
	return d
}

func (d DownloadError) GetFileSize() int64 {
	return 100
}

func (d DownloadError) GetFileName() string {
	return "data"
}

func (d DownloadError) Make(ctx context.Context) (DomainLevel.Download, error) {
	return d, nil
}

func (d DownloadError) SetAdditionalReader(reader io.Reader) DomainLevel.MakerDownloader {
	return d
}
