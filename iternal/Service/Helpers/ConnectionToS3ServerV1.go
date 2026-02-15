package Helpers

import (
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func Inzelire() (*session.Session, error) {

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("ru-1"),
		Endpoint:         aws.String("https://s3.twcstorage.ru"),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretKey,
			"",
		),
	})

	if err != nil {
		slog.Error("err", err)
		return nil, err
	}

	return sess, nil
}
