package s3Interation

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Variables struct {
	Bucket     string
	S3Connect  *s3.Client
	OldConnect *session.Session
}
