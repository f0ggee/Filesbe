package Handlers

import (
	"Kaban/iternal/InfrastructureLayer"
	Uttiltesss2 "Kaban/iternal/Service/Helpers"
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
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func DownloadEncrypt(w http.ResponseWriter, ctxs context.Context, name string) error {

	ctx, cancel := Uttiltesss2.ContextForDownloading(ctxs)
	defer cancel()
	apps := *InfrastructureLayer.ConnectKeyControl()

	RealNameFile := apps.Key.GetRealNameFile(name)
	aesKey, err := apps.Key.ProcessingFileParameters(RealNameFile)

	if err != nil {
		return err
	}
	//Return the string as bytes
	ReturnToBytes, err := hex.DecodeString(aesKey)
	if err != nil {
		slog.Error("Error decode to string", "Err", err)
		return err
	}
	Mut.RLock()
	PrivateKey := NewPrivateKey
	Mut.RUnlock()
	//Decrypt the key
	AesDecryptKey, err := rsa.DecryptOAEP(sha256.New(), nil, PrivateKey, ReturnToBytes, nil)

	//Processing errors
	switch {
	//if our key has been changed, we try to use the  old key
	case strings.Contains(fmt.Sprint(err), "decryption error"):
		Mut.RLock()
		oldKey := OldPrivateKey
		Mut.RUnlock()
		slog.Error("Key is old")
		AesDecryptKey, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, oldKey, ReturnToBytes, nil)
		if err != nil {
			slog.Error("Error also decrypt with an old key ", err)
			return err
		}

	case err != nil:
		slog.Error("Error decode string", "Err", err)
		slog.Info("Info name", name)
		return err
	}

	Reader, writer := io.Pipe()

	sees, err := Uttiltesss2.Inzelire()
	if err != nil {
		slog.Error("Error in func", err)

		return err
	}
	downloader := s3.New(sees)

	o, err := downloader.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket:      aws.String(Bucket),
		IfNoneMatch: aws.String(""),
		Key:         aws.String(RealNameFile),
	})

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Error closing body", "Err", err)
			return

		}
	}(o.Body)

	switch {
	case strings.Contains(fmt.Sprint(err), "NoSuchKey"):
		slog.Info("File was used")
		return errors.New("file was used")

	case errors.Is(err, context.DeadlineExceeded):
		slog.Error("Time was exceeded")
		return errors.New("time was exceeded")
	case errors.Is(err, context.Canceled):
		slog.Info("a user has been cancelled download ")
		return errors.New("a user has been canceled download ")
	case err != nil:
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
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename= %v", RealNameFile))
	w.Header().Set("Content-Length", strconv.FormatInt(*o.ContentLength-aes.BlockSize, 10))

	if _, err = io.Copy(w, Reader); err != nil {
		slog.Error("Err In file Service Downloader Encrypt", "err", err)
		return errors.New("connect close")
	}

	slog.Info("Func Download ends")
	//Here, we start the func, which deletes a file
	DeleteFile(RealNameFile, false)

	return nil
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
