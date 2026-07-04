package DomainLevel

import (
	"context"
	"io"
	"mime/multipart"

	NewVersion "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/service/s3"
)

type DeleterS3 interface {
	DeleteFileFromS3(string, context.Context) error
	DeleterS3Test(string, context.Context) error
}

type S3Uploader interface {
	UploadFile(parts int, goroutines int, ctx context.Context, fileFormat string, fileName string, file multipart.File) error

	UploadFileEncrypt(BesParts int, goroutine int, ctx context.Context, shortFileName string, ContentType string, reader *io.PipeReader) error
}

type DownloadingS3 interface {
	GetDownload(string, context.Context) (*NewVersion.GetObjectOutput, error)
	GetDownloadSecure(context.Context, string) (*s3.GetObjectOutput, error)
}
