package Uttiltesss

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"log/slog"
)

func Inzelire() (*session.Session, error) {

	acsee_key := "KGS5ANSMLR5IXDJPH85H"
	secret_key := "hTvM2Qg5HqDvHx2vHVbQePEjbmK8XWgsqukcwsmn"
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("ru-1"),
		Endpoint:         aws.String("https://s3.twcstorage.ru"),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			acsee_key,
			secret_key,
			"",
		),
	})

	if err != nil {
		slog.Error("err", err)
		return nil, err
	}

	return sess, nil
}
