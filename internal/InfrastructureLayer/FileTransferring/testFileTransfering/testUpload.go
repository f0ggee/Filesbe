package testFileTransfering

import (
	"Kaban/internal/DomainLevel"
	"context"
	"fmt"
	"io"
)

var Hub []byte

type FirstUploading struct {
}

func GetFirstUploading() FirstUploading {
	return FirstUploading{}
}
func (f *FirstUploading) CloseSource() error {
	return nil

}
func (f *FirstUploading) Upload(r io.Reader) error {
	buf := make([]byte, 1024)
	for {

		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		Hub = append(Hub, buf[:n]...)
	}
	return nil
}

func (f FirstUploading) SetName(s string) DomainLevel.MakerUploader {
	return f
}

func (f FirstUploading) SetSize(i int64) DomainLevel.MakerUploader {
	return f
}

func (f FirstUploading) GetFileSize() int64 {
	return 30
}

func (f FirstUploading) GetFileName() string {
	return "er"
}

func (f FirstUploading) Make(ctx context.Context) (DomainLevel.Upload, error) {
	return &f, nil
}

func (f FirstUploading) SetAdditionalWriter(writer io.Writer) DomainLevel.MakerUploader {
	return f
}

type UploaderSecond struct {
}

func GetUploaderSecond() UploaderSecond {
	return UploaderSecond{}
}

func (u *UploaderSecond) CloseSource() error {

	return nil

}

func (u *UploaderSecond) Upload(r io.Reader) error {

	buf := make([]byte, 1024)

	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fmt.Println("Test Uploader:", string(buf[:n]))
	}
	return nil
}

func (u *UploaderSecond) GetFileName() string {
	return "Some"
}
func (u *UploaderSecond) GetFileSize() int64 {
	return 1000
}

func (u *UploaderSecond) Make(context.Context) (DomainLevel.Upload, error) {
	return u, nil
}

func (u *UploaderSecond) SetAdditionalWriter(io.Writer) DomainLevel.MakerUploader {
	return u
}

func (u *UploaderSecond) SetName(string) DomainLevel.MakerUploader {
	return u
}

func (u *UploaderSecond) SetSize(int64) DomainLevel.MakerUploader {
	return u
}
