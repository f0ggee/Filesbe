package DomainLevel

import (
	"context"
	"io"

	NewVersion "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/service/s3"
)

type DeleterS3 interface {
	DeleteFileFromS3(string, context.Context) error
	DeleterS3Test(string, context.Context) error
}

type TypeUploading struct {
	Pipe   *io.PipeReader
	Normal io.ReadCloser
}

type FileDetails struct {
	FileFormat string
	FileName   string
	FileBody   TypeUploading
}
type UploadFileIncomingData struct {
	Parts      int
	Goroutines int
	Ctx        context.Context
	FileDetails
}

type S3Uploader interface {
	UploadFile(UploadFileIncomingData) error

	UploadFileEncrypt(UploadFileIncomingData) error
}

type DownloadingS3 interface {
	GetDownload(string, context.Context) (*NewVersion.GetObjectOutput, error)
	GetDownloadSecure(context.Context, string) (*s3.GetObjectOutput, error)
}
