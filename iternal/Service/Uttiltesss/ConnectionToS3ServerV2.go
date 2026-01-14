package Uttiltesss

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func Inzelire2() (*s3.Client, error) {
	acsee_key := os.Getenv("Access_Key")
	secret_key := os.Getenv("Secret_key")
	Endpoiont := os.Getenv("end")
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
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(acsee_key, secret_key, "")),
		config.WithRegion("ru-1"))
	if err != nil {
		slog.Error("Erro in Inzelire", err)
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(Endpoiont)

	})
	slog.Info("Connect to S3 server success")
	return client, nil
}
