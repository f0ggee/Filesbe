package DeleterS3

import "github.com/aws/aws-sdk-go-v2/service/s3"

type DeleterS3 struct {
	Conf *s3.Client
}
