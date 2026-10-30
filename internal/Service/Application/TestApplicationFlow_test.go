package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/FileTransferring/testFileTransfering"
	"testing"
)

type flow struct {
	FirstUpload  DomainLevel.MakerUploader
	Download     DomainLevel.MakerDownloader
	SecondUpload DomainLevel.MakerUploader
}

var flowControl flow

func makeFlow() flow {

	return flow{
		FirstUpload:  testFileTransfering.GetFirstUploading(),
		Download:     testFileTransfering.GetTestDownloader(),
		SecondUpload: testFileTransfering.GetUploaderSecond(),
	}
}
func TestApp_test(t *testing.T) {

	t.Run("Upload_Secure", func(t *testing.T) {
		TestUploadEncrypt_UploadEncryptApplication(t)
	})
	t.Run("Download_secure", func(t *testing.T) {

	})
}
