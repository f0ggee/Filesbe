package Handlers

import (
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
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

func (sa *HandlerPackCollect) DownloadEncrypt(w http.ResponseWriter, ctxs context.Context, name string) error {

	fileInfoInBytes, err := sa.RedisControlling.Reader.GetFileInfo(name, ctxs)
	if err != nil {
		return err
	}

	Keys.Mut.RLock()
	newPrivateKey := Keys.NewPrivateKey
	oldPrivateKey := Keys.OldPrivateKey
	Keys.Mut.RUnlock()
	aesKey, realFileName, err := sa.Crypto.Decrypt.DecryptFileInfo(fileInfoInBytes, newPrivateKey.Bytes(), oldPrivateKey.Bytes())
	if err != nil {
		return err
	}

	Reader, writer := io.Pipe()

	err = sa.downloadFileToClient(w, ctxs, name, writer, aesKey, realFileName, Reader)
	if err != nil {
		return err
	}

	slog.Info("Func Delete start in download encrypt")

	err = sa.RedisControlling.Deleter.DeleteFileInfo(name, ctxs)
	if err != nil {
		return err
	}
	err = sa.S3.Deleter.DeleteFileFromS3(name, ctxs)
	if err != nil {
		return err
	}

	slog.Info("Func Delete ends in download encrypt  ")

	return nil
}

func (sa *HandlerPackCollect) downloadFileToClient(w http.ResponseWriter, ctx context.Context, name string, writer *io.PipeWriter, aesKey []byte, realFileName string, Reader *io.PipeReader) error {

	downloader := s3.New(sa.S3.S3OldConnect)

	o, err := downloader.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket:      aws.String(Bucket),
		IfNoneMatch: aws.String(""),
		Key:         aws.String(name),
	})

	defer func(Body io.ReadCloser) {
		err = Body.Close()
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

	}
	if err != nil {
		slog.Error("ServiceDownload:", err.Error())
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
		err = DecryptFile(aesKey, o, writer)
		if err != nil {
			err := writer.CloseWithError(err)
			if err != nil {
				slog.Error("Error", "err", err)
				return
			}

		}
	}()

	FormatFile := sa.FileInfo.FileManaging.FindFormatOfFile(realFileName)
	w.Header().Set("Content-Type", FormatFile)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename= %v", realFileName))
	w.Header().Set("Content-Length", strconv.FormatInt(*o.ContentLength-aes.BlockSize, 10))

	if _, err = io.Copy(w, Reader); err != nil {
		slog.Error("Err In file Service Downloader EncryptFile", "err", err)
		return errors.New("connect close")
	}
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
		slog.Error("Error in  create file", err.Error())
		return err
	}

	nonce := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(o.Body, nonce)
	if err != nil {
		slog.Error("Error in read", err.Error())
		return err
	}

	plaintext := make([]byte, 35*1024)

	stream := cipher.NewCTR(block, nonce)

	file := bufio.NewReader(o.Body)
	for {
		n, err := file.Read(plaintext)
		if err != nil && err != io.EOF {
			slog.Error("Error in file", err.Error())
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
