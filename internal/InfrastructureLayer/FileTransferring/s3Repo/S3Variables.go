package s3Repo

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Variables struct {
	Bucket    string
	S3Connect *s3.Client
}

func GetNewVariables(bucket string, s3Connect *s3.Client) *Variables {
	return &Variables{Bucket: bucket, S3Connect: s3Connect}
}
