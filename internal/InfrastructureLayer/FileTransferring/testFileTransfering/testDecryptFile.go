package testFileTransfering

import (
	"Kaban/internal/DomainLevel"
	"bytes"
	"context"
	"errors"
	"io"
)

type TestDecryptFile struct {
	U          DomainLevel.MakerUploader
	D          DomainLevel.MakerDownloader
	h          *[]byte
	cipherText *[]byte
}

func (t *TestDecryptFile) SetCipherText(data []byte) {
	*t.cipherText = data
}

func (t *TestDecryptFile) Reset() {
	*t.h = nil
	*t.cipherText = nil
}

func GetNormalFlow() TestDecryptFile {
	var h []byte
	var _cipherText []byte

	return TestDecryptFile{
		U: testUploader{
			d: &h,
		},
		D: testDownload{
			h: &_cipherText,
		},
		h:          &h,
		cipherText: &_cipherText,
	}
}

type testDownloadFile struct {
	cipherText *[]byte
}

func (t *testDownloadFile) Downloader() (io.Reader, error) {
	return bytes.NewReader(*t.cipherText), nil
}

func (t *testDownloadFile) CloseSource() error {
	return nil
}

func (t *testDownloadFile) SetName(s string) DomainLevel.MakerDownloader {
	return t
}

func (t *testDownloadFile) SetSize(i int64) DomainLevel.MakerDownloader {
	return t
}

func (t *testDownloadFile) GetFileSize() int64 {
	return int64(len(*t.cipherText))
}

func (t *testDownloadFile) GetFileName() string {
	return "data"
}

func (t *testDownloadFile) Make(ctx context.Context) (DomainLevel.Download, error) {
	return t, nil
}

func (t *testDownloadFile) SetAdditionalReader(reader io.Reader) DomainLevel.MakerDownloader {
	return t
}

type testUploadFile struct {
	h *[]byte
}

func (t *testUploadFile) Upload(reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	*t.h = append(*t.h, data...)
	return nil
}

func (t *testUploadFile) CloseSource() error {
	return nil
}

func (t *testUploadFile) SetName(s string) DomainLevel.MakerUploader {
	return t
}

func (t *testUploadFile) SetSize(i int64) DomainLevel.MakerUploader {
	return t
}

func (t *testUploadFile) GetFileSize() int64 {
	return int64(len(*t.h))
}

func (t *testUploadFile) GetFileName() string {
	return "data"
}

func (t *testUploadFile) Make(ctx context.Context) (DomainLevel.Upload, error) {
	return t, nil
}

func (t *testUploadFile) SetAdditionalWriter(writer io.Writer) DomainLevel.MakerUploader {
	return t
}

func GetTestErrorFlowErrorDecryptFile() TestDecryptFile {

	var h []byte

	var cipherText []byte
	return TestDecryptFile{
		U: testUploadFileError{
			h: &h,
		},
		D: testDownloadFileError{
			cipherText: &cipherText,
		},
		h:          &h,
		cipherText: &cipherText,
	}
}

type testUploadFileError struct {
	h *[]byte
}

func (t testUploadFileError) Upload(reader io.Reader) error {
	return errors.New("error")
}

func (t testUploadFileError) CloseSource() error {
	return nil
}

func (t testUploadFileError) SetName(s string) DomainLevel.MakerUploader {
	return t
}

func (t testUploadFileError) SetSize(i int64) DomainLevel.MakerUploader {
	return t
}

func (t testUploadFileError) GetFileSize() int64 {
	if t.h == nil {
		return 0
	}
	return int64(len(*t.h))
}

func (t testUploadFileError) GetFileName() string {
	return "data"
}

func (t testUploadFileError) Make(ctx context.Context) (DomainLevel.Upload, error) {
	return t, nil
}

func (t testUploadFileError) SetAdditionalWriter(writer io.Writer) DomainLevel.MakerUploader {
	return t
}

type testDownloadFileError struct {
	cipherText *[]byte
}

func (t testDownloadFileError) Downloader() (io.Reader, error) {
	return nil, errors.New("error")
}

func (t testDownloadFileError) CloseSource() error {
	return nil
}

func (t testDownloadFileError) SetName(s string) DomainLevel.MakerDownloader {
	return t
}

func (t testDownloadFileError) SetSize(i int64) DomainLevel.MakerDownloader {
	return t
}

func (t testDownloadFileError) GetFileSize() int64 {
	if t.cipherText == nil {
		return 0
	}
	return int64(len(*t.cipherText))
}

func (t testDownloadFileError) GetFileName() string {
	return "data"
}

func (t testDownloadFileError) Make(ctx context.Context) (DomainLevel.Download, error) {
	return t, nil
}

func (t testDownloadFileError) SetAdditionalReader(reader io.Reader) DomainLevel.MakerDownloader {
	return t
}
