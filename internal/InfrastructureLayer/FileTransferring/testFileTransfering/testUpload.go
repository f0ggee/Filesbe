package testFileTransfering

import (
	"Kaban/internal/DomainLevel"
	"context"
	"io"
	"os"
)

type Uploader struct {
	f *os.File
}

func NewUploader() *Uploader {
	return &Uploader{}
}

func (u Uploader) Uploader(reader io.Reader) error {

	buf := make([]byte, 32)

	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		_, err = u.f.Write(buf[:n])
		if err != nil {
			panic(err)
		}

	}
	return nil
}

func (u Uploader) CloseSource() error {
	//TODO implement me
	panic("implement me")
}

func (u Uploader) SetName(s string) DomainLevel.MakerUploader {
	//TODO implement me
	panic("implement me")
}

func (u Uploader) SetSize(i int64) DomainLevel.MakerUploader {
	//TODO implement me
	panic("implement me")
}

func (u Uploader) GetFileSize() int64 {
	//TODO implement me
	panic("implement me")
}

func (u Uploader) GetFileName() string {
	//TODO implement me
	panic("implement me")
}

func (u *Uploader) Make(ctx context.Context) (DomainLevel.Upload, error) {

	file, err := os.Create("example1.txt")
	if err != nil {
		panic(err)
	}

	u.f = file
	return u, nil
}

func (u Uploader) SetAdditionalWriter(writer io.Writer) DomainLevel.MakerUploader {
	//TODO implement me
	panic("implement me")
}
