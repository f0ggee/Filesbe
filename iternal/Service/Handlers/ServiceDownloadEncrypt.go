package Handlers

import (
	Dto2 "Kaban/iternal/Dto"
	Uttiltesss2 "Kaban/iternal/Service/Uttiltesss"
	"bufio"
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
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func DownloadEncrypt(w http.ResponseWriter, ctxs context.Context, name string) error {

	ctx, cancel := Uttiltesss2.Contexte(ctxs)
	defer cancel()
	NameOfFile := name

	NameOfFile = GetTrueName(NameOfFile, name)

	sa, err2 := TreatOfFile(NameOfFile)
	if err2 != nil {
		return err2
	}

	defer func() {
		slog.Info("Func delete starts ")
		DeleteFile(name)
	}()

	ReturnToBytes, err := hex.DecodeString(sa.AesKey)
	if err != nil {
		slog.Error("Error decode to string", "Err", err)
		return err
	}
	AesDecryptKey, err := rsa.DecryptOAEP(sha256.New(), nil, NewPrivateKey, ReturnToBytes, nil)

	switch {
	case errors.Is(err, errors.New("decryption error")):
		AesDecryptKey, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, OldPrivateKey, ReturnToBytes, nil)
		if err != nil {
			slog.Error("Error also decrypt with an old key ", err)
			return err
		}

	}
	if err != nil {
		slog.Error("Error decode string", "Err", err)
		slog.Info("Info name", name)
		return err
	}
	Reader, writer := io.Pipe()

	bucket := os.Getenv("BUCKET")

	sees, err := Uttiltesss2.Inzelire()
	if err != nil {
		slog.Error("Error in func", err)

		return err
	}
	downloader := s3.New(sees)

	o, err := downloader.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket:      aws.String(bucket),
		IfNoneMatch: aws.String(""),
		Key:         aws.String(NameOfFile),
	})

	if err != nil {
		slog.Error("ServiceDownload:", err)
		return err
	}

	go func() {
		defer func(writer *io.PipeWriter) {
			err := writer.Close()
			if err != nil {
				slog.Error("Writer can't close", "Err", err)
				return
			}
		}(writer)
		err = DecryptFile(AesDecryptKey, o, writer)
		if err != nil {
			err := writer.CloseWithError(err)
			if err != nil {
				slog.Error("Error", "err", err)
				return
			}

		}
	}()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename= %v", NameOfFile))
	w.Header().Set("Content-Length", strconv.FormatInt(*o.ContentLength-aes.BlockSize, 10))

	if _, err = io.Copy(w, Reader); err != nil {
		slog.Error("Err In file Service Downloader Encrypt", "err", errors.New("connect close"))
		return errors.New("connect close")
	}

	return nil
}

func TreatOfFile(NameOfFile string) (struct {
	AesKey          string
	TimeSet         time.Time
	IsUsed          bool
	IsStartDownload bool
}, error) {
	sa, ok := Dto2.MapForFile[NameOfFile]

	if !ok {
		slog.Error("File don't find ")
		return struct {
			AesKey          string
			TimeSet         time.Time
			IsUsed          bool
			IsStartDownload bool
		}{}, errors.New("don't find ")
	}
	if sa.IsUsed {
		return struct {
			AesKey          string
			TimeSet         time.Time
			IsUsed          bool
			IsStartDownload bool
		}{}, errors.New("file was used")
	}
	sa.IsStartDownload = true
	return sa, nil
}

// If i want to save files
func treateOfFile(NameOfFile string) error {
	sa, ok := Dto2.MapForFile[NameOfFile]

	if !ok {
		slog.Error("File don't find ")
		return errors.New("don't find ")
	}
	if sa.IsUsed {
		return errors.New("file was used")
	}
	sa.IsStartDownload = true
	return nil
}

func GetTrueName(NameOfFile string, name string) string {
	if Name, ok := Dto2.NamesToConvert[NameOfFile]; ok {
		NameOfFile = Name
		delete(Dto2.NamesToConvert, name)
	}
	return NameOfFile
}

func DecryptFile(AesKey []byte, o *s3.GetObjectOutput, writer *io.PipeWriter) error {

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Error close s3 connection ", "err", err)
			return

		}
	}(o.Body)

	block, err := aes.NewCipher(AesKey)
	if err != nil {
		slog.Error("Error in  create file", err)
		return err
	}

	nonce := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(o.Body, nonce)
	if err != nil {
		slog.Error("Error in read", err)
		return err
	}

	plaintext := make([]byte, 35*1024)

	stream := cipher.NewCTR(block, nonce)

	file := bufio.NewReader(o.Body)
	for {
		n, err := file.Read(plaintext)
		if err != nil && err != io.EOF {
			slog.Error("Error in file", err)
			return err
		}
		if err == io.EOF {
			break
		}

		if n > 0 {
			stream.XORKeyStream(plaintext[:n], plaintext[:n])
			_, err = writer.Write(plaintext[:n])
			if err != nil {
				err := writer.CloseWithError(err)
				if err != nil {
					slog.Error("Error is writing into file", "Err", err)
					return err
				}
				return err
			}
		}

	}

	return nil
}
