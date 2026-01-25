package Handlers

import (
	"Kaban/iternal/Dto"
	"Kaban/iternal/Service/Helpers"
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func DeleteFile(Name string, IsEncrypt bool) {
	slog.Info("Func deleteFile starts")

	nameOfFile := Name

	if IsEncrypt {
		if s, _ := Dto.MapForFile[nameOfFile]; s.IsStartDownload == true {
			slog.Info("Flag is true")
			return
		}
	}

	cfgs, err := Helpers.S3Helper()
	if err != nil {
		slog.Error("can't connect to S3 server", "Err", err)
		return
	}
	s := &s3.DeleteObjectInput{
		Bucket: aws.String(Bucket),
		Key:    &nameOfFile,
	}

	_, err = cfgs.DeleteObject(context.Background(), s)
	if err != nil {
		slog.Error("Error in delete func", err)

	}

	slog.Info("Func deleteFile ends")
	return

}
