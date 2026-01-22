package Handlers

import (
	"Kaban/iternal/Dto"
	Uttiltesss2 "Kaban/iternal/Service/Helpers"
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

	reader, writer := io.Pipe()

	file, sizeAndName, err := r.FormFile("file")
	if err != nil {
		slog.Error("Err from FileUploader 1 ", err)
		http.Error(w, "Error", http.StatusNotFound)
		return "", errors.New("can't get file")
	}

	//This function checks len of name
	NameS := Uttiltesss2.CheckLenOfName(sizeAndName.Filename)

	// The function below  finds  best options for download
	SizeFile := sizeAndName.Size
	BesParts, goroutine := Uttiltesss2.FindBest(SizeFile)

	defer func() {
		err = file.Close()
		if err != nil {
			slog.Error("Err, cant' close a file", "err", err)
			return
		}
	}()

	ctx, cancel := Uttiltesss2.ContextForDownloading(r.Context())
	if err != nil {
		slog.Error("Error in Check File size", err)
		return "", errors.New("file size too big")
	}
	defer cancel()
	timeS := time.Now()

	defer func() {
		sa := time.Since(timeS)
		fmt.Println(sa)
	}()
	cfg, err := Uttiltesss2.Initialization2()
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

	publicKey := &NewPrivateKey.PublicKey

	encryptAesKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, <-chanelForAesKey, nil)
	if err != nil {
		slog.Error("Error create aes", "ERR", err)
		return "", errors.New("can't validate data ")
	}

	IntoString := hex.EncodeToString(encryptAesKey)

	Dto.MapForFile[sizeAndName.Filename] = struct {
		AesKey          string
		TimeSet         time.Time
		IsUsed          bool
		IsStartDownload bool
	}{AesKey: IntoString, TimeSet: time.Now(), IsUsed: false, IsStartDownload: false}

	fmt.Println(BesParts, goroutine)
	uploader := manager.NewUploader(cfg, func(uploader *manager.Uploader) {
		uploader.MaxUploadParts = 200
		uploader.PartSize = int64(BesParts) * 1024 * 1024
		uploader.Concurrency = goroutine
	})

	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(Bucket),
		Key:    aws.String(sizeAndName.Filename),
		Body:   reader,
	})

	var ns *types.NoSuchKey
	switch {
	case errors.As(err, &ns):

		slog.Error("file was used")
		return "", errors.New("file was used")

	}
	if err != nil {
		slog.Error("Error in uploader", err)
		return "", errors.New("error in uploader file")
	}

	//The Func  deletes a file
	time.AfterFunc(5*time.Minute, func() {

		DeleteFile(NameS, true)

	})

	slog.Info("File success upload ")

	return NameS, nil

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
