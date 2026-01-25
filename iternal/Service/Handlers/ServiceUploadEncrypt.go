package Handlers

import (
	"Kaban/iternal/InfrastructureLayer"
	Uttiltesss2 "Kaban/iternal/Service/Helpers"
	"Kaban/iternal/Service/Helpers/validator"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
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

func FileUploaderEncrypt(w http.ResponseWriter, r *http.Request) (string, error) {
	slog.Info("Func FileUploaderEncrypt starts")
	file, sizeAndName, err := r.FormFile("file")
	if err != nil {
		slog.Error("Err from FileUploader 1 ", err)
		http.Error(w, "Error", http.StatusNotFound)
		return "", errors.New("can't get file")
	}

	defer func() {
		err = file.Close()
		if err != nil {
			slog.Error("Err, cant' close a file", "err", err)
			return
		}
	}()
	apps := *InfrastructureLayer.ConnectKeyControl()

	err = validator.CheckFileSize2(sizeAndName.Size)
	if err != nil {
		slog.Error("check file size error", err)
		return "", err
	}

	ctx, cancel := Uttiltesss2.Context2(r.Context())
	if cancel == nil {
		return "", errors.New("error in file Uploader no encrypt")
	}
	defer cancel()

	reader, writer := io.Pipe()

	//This function checks len of name
	NameS := CheckLenOfName(sizeAndName.Filename)

	// The function below  finds  best options for download
	BesParts, goroutine := Uttiltesss2.FindBest(sizeAndName.Size)

	timeS := time.Now()

	defer func() {
		sa := time.Since(timeS)
		fmt.Println(sa)
	}()
	cfg, err := Uttiltesss2.S3Helper()
	if err != nil {
		slog.Error("can't connect to S3 server", "Err", err)
		return "", errors.New("can't connect to our servers")
	}

	chanelForAesKey := make(chan []byte, 100)

	// The function, which encrypts a file
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

	AesKeyIntoString, err2 := encryptKey(chanelForAesKey)
	if err2 != nil {
		return "", err2
	}

	apps.Key.SaveFileInfo(sizeAndName.Filename, AesKeyIntoString)

	_, err3 := uploadFileEncrypt(cfg, BesParts, goroutine, ctx, sizeAndName, reader)
	if err3 != nil {
		return "", err3
	}

	//The Func  deletes a file
	time.AfterFunc(5*time.Minute, func() {

		DeleteFile(NameS, true)

	})

	slog.Info("File success upload ")

	return NameS, nil

}

func uploadFileEncrypt(cfg *s3.Client, BesParts int, goroutine int, ctx context.Context, sizeAndName *multipart.FileHeader, reader *io.PipeReader) (string, error) {
	uploader := manager.NewUploader(cfg, func(uploader *manager.Uploader) {
		uploader.MaxUploadParts = 200
		uploader.PartSize = int64(BesParts) * 1024 * 1024
		uploader.Concurrency = goroutine
	})

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(Bucket),
		Key:    aws.String(sizeAndName.Filename),
		Body:   reader,
	})
	if err == nil {
		return "", nil
	}

	var ns *types.NoSuchKey

	switch {

	case errors.As(err, &ns):

		slog.Error("file was used")
		return "", errors.New("file was used")

	case errors.Is(err, context.Canceled):
		slog.Error("file was cancelled")
		return "", errors.New("file was cancelled")

	case errors.Is(err, context.DeadlineExceeded):
		slog.Error("Time was exceeded")
		return "", errors.New("time was exceeded")

	default:
		slog.Error("file was not uploadable", err.Error())
		return "", errors.New("file was not uploadable")
	}
}

func encryptKey(chanelForAesKey chan []byte) (string, error) {

	for {

		key, _ := <-chanelForAesKey

		publicKey := &NewPrivateKey.PublicKey

		encryptAesKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, key, nil)
		if err != nil {
			slog.Error("Error create aes", "ERR", err)
			return "", errors.New("can't validate data ")
		}

		IntoString := hex.EncodeToString(encryptAesKey)
		return IntoString, nil

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
			slog.Error("Error in file upload", err)
			return err
		}
		if err == io.EOF {
			break
		}
		stream.XORKeyStream(buf[:n], buf[:n])
		_, err = writer.Write(buf[:n])
		if err != nil {
			slog.Error("Err write in process", err)
			return err
		}

	}

	return nil
}
