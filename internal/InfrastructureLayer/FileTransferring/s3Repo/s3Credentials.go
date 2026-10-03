package s3Repo

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var accessKey = os.Getenv("Access_Key")
var secretKey = os.Getenv("Secret_key")
var EndPoint = os.Getenv("end")

var s3Cred *ConnectCredentials

func init() {

	s, err := GenerateS3()
	if err != nil {
		panic(err)
	}
	s3Cred = s
}

type ConnectCredentials struct {
	bucket string
	conn   *s3.Client
}

func (receiver ConnectCredentials) getConnect() *s3.Client {
	return receiver.conn
}

func (receiver ConnectCredentials) getBucket() string {
	return receiver.bucket
}

func GenerateS3() (*ConnectCredentials, error) {
	s, err := EstablishS3()
	if err != nil {
		return nil, err
	}
	buck := os.Getenv("BUCKET")
	if buck == "" {
		slog.Error("GenerateS3: a bucket is empty")
		return nil, errors.New(ErrorBucket)
	}
	return &ConnectCredentials{
		bucket: buck,
		conn:   s,
	}, nil
}

const ErrorS3Credentials = "the key wasn't created"
const ErrorConnectS3 = "can't connect to a s3 server"
const ErrorBucket = "a bucket wasn't created"

func EstablishS3() (*s3.Client, error) {
	tr := &http.Transport{
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 300,
		MaxConnsPerHost:     300,
	}
	if accessKey == "" || secretKey == "" {
		slog.Error("EstablishS3: error either a secret key or an access key wasn't created")
		return nil, errors.New(ErrorS3Credentials)
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithHTTPClient(&http.Client{Transport: tr}),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRetryMaxAttempts(6),
		config.WithRetryMode("adaptive"),
		config.WithRegion("ru-1"))
	config.WithBaseEndpoint(EndPoint)
	if err != nil {
		slog.Error("EstablishS3:error to establish connect", "ERROR", err)
		return nil, errors.New(ErrorConnectS3)
	}
	client := s3.NewFromConfig(cfg, func(options *s3.Options) {})
	return client, nil
}
