package DomainLevel

import (
	"context"
	"io"
)

type Upload interface {
	Uploader(io.Reader) error
	CloseSource() error
}
type FileDetails interface {
	GetFileSize() int64
	GetFileName() string
}
type MakerUploader interface {
	SetName(string) MakerUploader
	SetSize(int64) MakerUploader
	FileDetails
	Make(context.Context) (Upload, error)
	SetAddWriter(io.Writer) MakerUploader
}
type Download interface {
	Downloader() (io.Reader, error)
	CloseSource() error
}
type MakerDownloader interface {
	SetName(string) MakerDownloader
	SetSize(int64) MakerDownloader
	FileDetails
	Make(ctx context.Context) (Download, error)
	SetAddReader(io.Reader) FileDetails
}

type Delete interface {
	Deleter() error
}

type MakerDeleter interface {
	SetName(string) MakerDeleter
	SetSize(int64) MakerDeleter

	Make(context.Context) (Delete, error)
}
