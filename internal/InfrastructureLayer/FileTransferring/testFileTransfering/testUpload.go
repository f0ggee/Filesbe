package testFileTransfering

import (
	"Kaban/internal/DomainLevel"
	"context"
	"io"
)

type hub struct {
	data []byte
}

type testUploader struct {
	d *hub
}

func (t testUploader) Upload(reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	t.d.data = append(t.d.data, data...)
	return nil
}

func (t testUploader) CloseSource() error {
	return nil
}

func (t testUploader) SetName(s string) DomainLevel.MakerUploader {
	return t
}

func (t testUploader) SetSize(i int64) DomainLevel.MakerUploader {
	return t
}

func (t testUploader) GetFileSize() int64 {
	return 100
}

func (t testUploader) GetFileName() string {
	return "data"
}

func (t testUploader) Make(ctx context.Context) (DomainLevel.Upload, error) {
	return t, nil
}

func (t testUploader) SetAdditionalWriter(writer io.Writer) DomainLevel.MakerUploader {
	return t
}

type ErrorUpload struct {
	err error
	ctx context.Context
}

func NewErrorUpload(err error) ErrorUpload {
	return ErrorUpload{err: err}
}

func (e ErrorUpload) Upload(reader io.Reader) error {
	if e.ctx != nil && e.ctx.Err() != nil {
		return e.ctx.Err()
	}
	return e.err
}

func (e ErrorUpload) CloseSource() error {
	return nil
}

func (e ErrorUpload) SetName(s string) DomainLevel.MakerUploader {
	return e
}

func (e ErrorUpload) SetSize(i int64) DomainLevel.MakerUploader {
	return e
}

func (e ErrorUpload) GetFileSize() int64 {
	return 100
}

func (e ErrorUpload) GetFileName() string {
	return "data"
}

func (e ErrorUpload) Make(ctx context.Context) (DomainLevel.Upload, error) {
	e.ctx = ctx
	return e, nil
}

func (e ErrorUpload) SetAdditionalWriter(writer io.Writer) DomainLevel.MakerUploader {
	return e
}
