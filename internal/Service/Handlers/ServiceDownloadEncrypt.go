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

	"golang.org/x/sync/errgroup"
)

func (sa *HandlerPackCollect) DownloadEncrypt(w http.ResponseWriter, ctxs context.Context, name string) error {

	fileInfoInBytes, err := sa.RedisControlling.Reader.GetFileInfo(name, ctxs)
	if err != nil {
		return err
	}

	newPrivateKey := sa.Keys.ControllerKey.GetKey()
	oldPrivateKey := sa.Keys.ControllerKey.GetOldKey()
	aesKey, realFileName, err := sa.Crypto.Decrypt.DecryptFileInfo(fileInfoInBytes, newPrivateKey, oldPrivateKey)
	if err != nil {
		return err
	}

	Reader, writer := io.Pipe()

	g, ctx := errgroup.WithContext(ctxs)

	Body, FileLength, err := sa.S3.S3Download.DownloadSecure(ctxs, name)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Error closing body", "error", err)
			return
		}
		return
	}(Body)

	g.Go(func() error {
		defer func(writer *io.PipeWriter) {
			err := writer.Close()
			if err != nil {
				slog.Error("Writer can't close", "Err", err)
				return
			}
		}(writer)
		err = DecryptFile(aesKey, Body, writer, ctx)
		if err != nil {
			err := writer.CloseWithError(err)
			if err != nil {
				slog.Error("Error", "err", err)
				return nil
			}

		}
		return nil
	})

	FormatFile := sa.FileInfo.FileManaging.FindFormatOfFile(realFileName)
	w.Header().Set("Content-Type", FormatFile)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename= %v", realFileName))
	w.Header().Set("Content-Length", strconv.FormatInt(FileLength-aes.BlockSize, 10))

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctxs.Err()
		default:
		}
		if _, err = io.Copy(w, Reader); err != nil {
			slog.Error("Err In file Service Downloader EncryptFile", "err", err)
			return errors.New("connect close")
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	slog.Info("Start deleting file")
	err = sa.RedisControlling.Deleter.DeleteFileInfo(name, ctxs)
	if err != nil {
		return err
	}
	err = sa.S3.Deleter.DeleteFileFromS3(name, ctxs)
	if err != nil {
		return err
	}
	slog.Info("Deleted file from S3")
	return nil
}

//func (sa *HandlerPackCollect) downloadFileToClient(w http.ResponseWriter, ctx context.Context, name string, writer *io.PipeWriter, aesKey []byte, realFileName string, Reader *io.PipeReader) error {
//	OldS3Connect, err := Helpers.Inzelire()
//	if err != nil {
//		slog.Error("Error connect to s3 old ", "Error", err)
//		return
//	}
//	downloader := *s3.New(OldS3Connect)
//
//	o, err := downloader.GetObjectWithContext(ctx, &s3.GetObjectInput{
//		Bucket:      aws.String(Bucket),
//		IfNoneMatch: aws.String(""),
//		Key:         aws.String(name),
//	})
//
//	defer func(Body io.ReadCloser) {
//		err = Body.Close()
//		if err != nil {
//			slog.Error("Error closing body", "Err", err)
//			return
//
//		}
//	}(o.Body)
//
//	switch {
//	case strings.Contains(fmt.Sprint(err), "NoSuchKey"):
//		slog.Info("File was used")
//		return errors.New("file was used")
//
//	case errors.Is(err, context.DeadlineExceeded):
//		slog.Error("Time was exceeded")
//		return errors.New("time was exceeded")
//	case errors.Is(err, context.Canceled):
//		slog.Info("a user has been cancelled download ")
//		return errors.New("a user has been canceled download ")
//
//	}
//	if err != nil {
//		slog.Error("ServiceDownload:", "Error", err.Error())
//		return err
//	}
//
//	g, ctx := errgroup.WithContext(ctx)
//	g.Go(func() error {
//		defer func(writer *io.PipeWriter) {
//			err := writer.Close()
//			if err != nil {
//				slog.Error("Writer can't close", "Err", err)
//				return
//			}
//		}(writer)
//		err = DecryptFile(aesKey, o.Body, writer, ctx)
//		if err != nil {
//			err := writer.CloseWithError(err)
//			if err != nil {
//				slog.Error("Error", "err", err)
//				return nil
//			}
//
//		}
//		return nil
//	})
//
//	return nil
//}

func DecryptFile(AesKey []byte, o io.ReadCloser, writer *io.PipeWriter, ctx context.Context) error {

	block, err := aes.NewCipher(AesKey)
	if err != nil {
		slog.Error("Error in  create file", "Error", err.Error())
		return err
	}

	nonce := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(o, nonce)
	if err != nil {
		slog.Error("Error in read", "Error", err.Error())
		return err
	}

	plaintext := make([]byte, 35*1024)

	stream := cipher.NewCTR(block, nonce)

	file := bufio.NewReader(o)
	for {
		if ctx.Err() != nil {
			return errors.New("context canceled")
		}
		n, err := file.Read(plaintext)
		if err != nil && err != io.EOF {
			slog.Error("Error in file", "Error", err.Error())
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
