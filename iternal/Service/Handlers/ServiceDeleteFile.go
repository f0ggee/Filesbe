package Handlers

import (
	"Kaban/iternal/Dto"
	"Kaban/iternal/Service/Uttiltesss"
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func DeleteFile(Name string) {

	nameOfFile := Name

	if s, _ := Dto.MapForFile[nameOfFile]; s.IsStartDownload == true {
		return
	}

	cfgs, err := Uttiltesss.Inzelire2()
	if err != nil {
		slog.Error("can't connect to S3 server", "Err", err)
		return
	}

	s := &s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET")),
		Key:    &nameOfFile,
	}

	_, err = cfgs.DeleteObject(context.Background(), s)
	if err != nil {
		slog.Error("Error in delete func", err)
		return

	}

	slog.Info("", slog.Group("Info about process"),
		slog.String("Name of file", nameOfFile),
		slog.Bool("Status", true))
}
