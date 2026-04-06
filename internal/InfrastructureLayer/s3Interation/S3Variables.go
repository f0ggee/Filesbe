package s3Interation

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Variables struct {
	Bucket    string
	S3Connect *s3.Client
}
