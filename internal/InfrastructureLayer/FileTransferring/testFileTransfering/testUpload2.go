package testFileTransfering

import (
	"Kaban/internal/DomainLevel"
	"context"
	"io"
	"os"
)

type Upload2 struct {
	f *os.File
}

func NewUpload2() *Upload2 {
	return &Upload2{}
}

func (u Upload2) Uploader(reader io.Reader) error {
	buf := make([]byte, 32)

	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		u.f.Write((buf[:n]))
	}
	return nil
}

func (u Upload2) CloseSource() error {
	//TODO implement me
	panic("implement me")
}

func (u *Upload2) SetName(s string) DomainLevel.MakerUploader {

	return u
}

func (u Upload2) SetSize(i int64) DomainLevel.MakerUploader {
	//TODO implement me
	panic("implement me")
}

func (u Upload2) GetFileSize() int64 {
	//TODO implement me
	panic("implement me")
}

func (u Upload2) GetFileName() string {
	//TODO implement me
	panic("implement me")
}

func (u *Upload2) Make(ctx context.Context) (DomainLevel.Upload, error) {
	f, err := os.Create("example2.txt")
	if err != nil {
		panic(err)
	}
	u.f = f
	return u, nil
}

func (u Upload2) SetAdditionalWriter(writer io.Writer) DomainLevel.MakerUploader {
	//TODO implement me
	panic("implement me")
}
