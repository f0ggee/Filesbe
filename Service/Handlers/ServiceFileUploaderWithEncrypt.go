package Handlers

import (
	"Kaban/Dto"
	"Kaban/Service/Uttiltesss"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func FileUploaderEncrypt(w http.ResponseWriter, r *http.Request) (string, error) {

	reader, writer := io.Pipe()

	file, sizeAndName, err := r.FormFile("file")
	if err != nil {
		slog.Error("Err from FileUploader 1 ", err)
		http.Error(w, "Error", http.StatusNotFound)
		return "", errors.New("can't get file")
	}

	NameS := sizeAndName.Filename

	// The function below  finds  best options for download

	SizeFile := sizeAndName.Size
	BesParts, Groutinse := FindBest(SizeFile)

	defer func() {
		err = file.Close()
		if err != nil {
			slog.Error("Err, cant' close a file", "err", err)
			return
		}
	}()

	ctx, cancel := Uttiltesss.Contexte(r.Context())
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
	cfg, err := Uttiltesss.Inzelire2()
	if err != nil {
		slog.Error("can't connect to S3 server", "Err", err)
		return "", errors.New("can't connect to our servers")
	}

	chanelForAesKey := make(chan []byte, 100)

	// The function, which encrypt a file
	go func() {
		defer func(writer *io.PipeWriter) {
			err := writer.Close()
			if err != nil {

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

	Dto.MapForFile[NameS] = struct {
		AesKey          string
		TimeSet         time.Time
		IsUsed          bool
		IsStartDownload bool
	}{AesKey: IntoString, TimeSet: time.Now(), IsUsed: false, IsStartDownload: false}

	fmt.Println(BesParts, Groutinse)
	uploader := manager.NewUploader(cfg, func(uploader *manager.Uploader) {
		uploader.MaxUploadParts = 200
		uploader.PartSize = int64(BesParts) * 1024 * 1024
		uploader.Concurrency = Groutinse
	})

	slog.Info("File", NameS)
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET")),
		Key:    aws.String(NameS),
		Body:   reader,
	})
	if err != nil {
		slog.Error("Error in uploader", err)
		return "", errors.New("error in uploader file")
	}

	//Func which in the end, deletes a file
	time.AfterFunc(5*time.Minute, func() {

		DeleteFile(NameS)

	})

	slog.Info("File success upload :)")

	return sizeAndName.Filename, nil

}

func Encrypt(file multipart.File, writer io.Writer, channelForBytes chan []byte) error {
	aesKey := make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		slog.Error("err generate aes-key", "Err", err)
		return errors.New("can't do advance protection")
	}

	slog.Info("Key", aesKey)

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

func FindBest(size int64) (int, int) {

	switch {
	case size > 500000000:

		fileResult := size / 1000000

		x := 50
		ResultPart := int(fileResult) / x
		for ResultPart > 15 {
			ResultPart = int(fileResult) / x
			x += 5

		}
		NumOfGroutine := ResultPart + 1
		return x, NumOfGroutine

	default:
		return 5, 20

	}
}
