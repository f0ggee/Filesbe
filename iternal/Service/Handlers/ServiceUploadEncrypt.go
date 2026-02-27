package Handlers

import (
	Uttiltesss2 "Kaban/iternal/Service/Helpers"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
)

const FileMaxSize = 1 << 30

func (sa *HandlerPackCollect) FileUploaderEncrypt(w http.ResponseWriter, r *http.Request) (string, error) {

	slog.Info("Func FileUploaderEncrypt starts")
	file, sizeAndName, err := r.FormFile("file")
	if err != nil {
		slog.Error("Err from FileUploader 1 ", err)
		return "", errors.New("can't get file")
	}
	if sizeAndName.Size >= FileMaxSize {
		slog.Info("File too big")

		return "", errors.New("file too big")
	}

	defer func() {
		err = file.Close()
		if err != nil {
			slog.Error("Err, cant' close a file", "err", err)
			return
		}
	}()

	ctx, cancel := Uttiltesss2.Context2(r.Context())
	if cancel == nil {
		return "", errors.New("error in file Uploader no encrypt")
	}
	defer cancel()

	reader, writer := io.Pipe()

	BesParts, goroutine := Uttiltesss2.FindBest(sizeAndName.Size)

	timeS := time.Now()

	defer func() {
		sa := time.Since(timeS)
		fmt.Println(sa)
	}()

	chanelForAesKey := make(chan []byte, 100)

	go func() {
		defer func(writer *io.PipeWriter) {
			err := writer.Close()
			if err != nil {
				slog.Error("can't close a file", "err", err)
				return
			}
		}(writer)
		err = Encrypt(file, writer, chanelForAesKey)
		if err != nil {
			err := writer.CloseWithError(err)
			if err != nil {
				slog.Error("Error in file writing", err)
				return
			}

		}
	}()

	AesKey := encryptKey(chanelForAesKey)

	FileInfoInBytes, err := sa.S.FileInfo.ConvertToBytesFileInfo(sizeAndName.Filename, AesKey)
	if err != nil {
		err := writer.CloseWithError(err)
		slog.Error("Error in file writing", err)
		return "", err
	}

	shortNameForFile := sa.S.FileInfo.GenerateShortFileName()
	//shortNameForFile := apps.Key.GenerateShortFileName()

	Keys.Mut.RLock()
	newPrivateKey := Keys.NewPrivateKey.Data()
	Keys.Mut.RUnlock()

	EncryptFileInfo, err := sa.S.FileInfo.EncryptData(FileInfoInBytes, newPrivateKey)
	if err != nil {
		slog.Error("Error in file writing", err)
		return "", err
	}

	FileExtension := sa.S.FileInfo.FindFormatOfFile(sizeAndName.Filename)
	_, err3 := uploadFileEncrypt(sa.S.S3Connect, BesParts, goroutine, ctx, shortNameForFile, FileExtension, reader)
	if err3 != nil {
		return "", err3
	}

	err = sa.S.RedisConn.WriteData(shortNameForFile, EncryptFileInfo)
	if err != nil {
		err := writer.CloseWithError(err)
		slog.Error("Error in file writing", err)
		return "", err
	}

	time.AfterFunc(5*time.Minute, func() {
		slog.Info("Func  Auto-FileDelete start")

		DownloadingHaveStarted := sa.S.RedisConn.ChekIsStartDownload(shortNameForFile)
		if DownloadingHaveStarted {
			return
		}
		err := sa.S.RedisConn.DeleteFileInfo(shortNameForFile)
		if err != nil {
			return
		}
		err = sa.S.S3Conn.DeleteFileFromS3(shortNameForFile, Bucket)
		if err != nil {
			return
		}
		slog.Info("Func Auto-deleteFile ends")

	})
	slog.Info("File success upload ")

	return shortNameForFile, nil

}

func uploadFileEncrypt(cfg *s3.Client, BesParts int, goroutine int, ctx context.Context, shortFileName string, ContentType string, reader *io.PipeReader) (string, error) {

	uploader := manager.NewUploader(cfg, func(uploader *manager.Uploader) {

		uploader.MaxUploadParts = 200
		uploader.PartSize = int64(BesParts) * 1024 * 1024
		uploader.Concurrency = goroutine
	})

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(Bucket),
		Key:         aws.String(shortFileName),
		ContentType: aws.String(ContentType),
		Body:        reader,
	})
	//Uploads, err := cfg.CreateMultipartUpload(context.Background(), &s3.CreateMultipartUploadInput{
	//	Bucket:      aws.String(Bucket),
	//	Key:         aws.String(shortFileName),
	//	ContentType: aws.String(ContentType),
	//}, func(options *s3.Options) {
	//
	//})

	if err == nil {
		return "", nil
	}

	var ns *types.NoSuchKey

	switch {

	case errors.As(err, &ns):

		slog.Error("file was")
		return "", errors.New("file was used")

	case errors.Is(err, context.Canceled):
		slog.Error("file downloading was cancelled")
		return "", errors.New("file was cancelled")

	case errors.Is(err, context.DeadlineExceeded):
		slog.Error("Time was exceeded")
		return "", errors.New("time was exceeded")

	default:
		slog.Error("file was not uploadable", err.Error())
		return "", errors.New("file was not uploadable")
	}
}

func encryptKey(chanelForAesKey chan []byte) []byte {

	for {

		key, _ := <-chanelForAesKey

		return key

	}
}

func Encrypt(file multipart.File, writer io.Writer, channelForBytes chan []byte) error {
	aesKey := make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		slog.Error("err generate aes-key", "Err", err)
		return errors.New("can't do advance protection")
	}

	channelForBytes <- aesKey
	close(channelForBytes)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		slog.Error("err create a NewCipher")
		return err
	}

	nonce := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	stream := cipher.NewCTR(block, nonce)
	buf := make([]byte, 32*1024)
	_, err = writer.Write(nonce)
	if err != nil {
		slog.Error("err write ", "Err", err)
		return err
	}
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			slog.Error("Error in file upload", err.Error())
			return err
		}
		if err == io.EOF {
			break
		}
		stream.XORKeyStream(buf[:n], buf[:n])
		_, err = writer.Write(buf[:n])
		if err != nil {
			slog.Error("Err write in process", err.Error())
			return err
		}

	}

	return nil
}
