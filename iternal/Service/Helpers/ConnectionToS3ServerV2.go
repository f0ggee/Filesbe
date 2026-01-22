package Helpers

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

var accessKey = os.Getenv("Access_Key")
var secretKey = os.Getenv("Secret_key")
var EndPoint = os.Getenv("end")

func init() {

	err := godotenv.Load("iternal/Service/.env")
	if err != nil {
		slog.Error("cannot load env file", err)
		return

	}

	accessKey = os.Getenv("Access_Key")
	secretKey = os.Getenv("Secret_key")
	EndPoint = os.Getenv("end")
}

func Initialization2() (*s3.Client, error) {

	//accessKey := os.Getenv("Access_Key")
	//secretKey := os.Getenv("Secret_key")
	//EndPoint := os.Getenv("end")

	tr := &http.Transport{
		MaxConnsPerHost:     300,
		MaxIdleConns:        512,
		MaxIdleConnsPerHost: 300,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  true,
	}
	httpClient := &http.Client{Transport: tr}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithHTTPClient(httpClient),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("ru-1"))
	if err != nil {
		slog.Error("Error Initialization", err)
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(EndPoint)

	})
	slog.Info("Connect to S3 server success")
	return client, nil
}
