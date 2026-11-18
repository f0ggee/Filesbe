package testFunctions

import (
	"Kaban/internal/Drivers/FileTransferring/TransferringHttpRepo"
	"Kaban/internal/Drivers/FileTransferring/s3Repo"
	"Kaban/internal/Drivers/FileTransferring/testFileTransfering"
	"Kaban/internal/Drivers/Parsers"
	"Kaban/internal/Drivers/RedisInteration"
	"Kaban/internal/Service/Application"
	"errors"
	"io"
)

func GetDownloadTest(testData io.Reader) *Application.NewDownload {
	dow := testFileTransfering.GetTestFlow()
	_u, err := dow.U.Make(nil)
	if err != nil {
		panic(err)
	}
	_u.Upload(testData)

	return &Application.NewDownload{
		DownloadFileControl: Application.DownloadFileControl{
			Downloader: dow.D,
			Uploader:   TransferringHttpRepo.NewHttpUploader(),
			Deleter:    s3Repo.NewS3DeleteTest(),
			Decoder:    Parsers.GetNewParsing(),
		},
		Reader: RedisInteration.NewRedisReadTest(),
	}

}
func GetDownloadErrorTest() *Application.NewDownload {
	return &Application.NewDownload{
		DownloadFileControl: Application.DownloadFileControl{
			Downloader: testFileTransfering.NewDownloadError(errors.New("error")),
			Uploader:   testFileTransfering.NewErrorUpload(errors.New("error")),
			Deleter:    nil,
			Decoder:    Parsers.GetNewParsing(),
		},
		Reader: RedisInteration.NewRedisReadTest(),
	}
}
