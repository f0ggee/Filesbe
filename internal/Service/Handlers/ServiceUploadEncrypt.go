package Handlers

import "C"
import (
	"Kaban/internal/Dto"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/awnumar/memguard"
	"golang.org/x/sync/errgroup"
)

func (sa *HandlerPackCollect) FileUploadEncrypt(r *http.Request) (string, error) {

	slog.Info("Func FileUploadEncrypt starts")
	file, sizeAndName, err := r.FormFile("file")
	if err != nil {
		slog.Error("Err from FileUploader 1 ", "Error", err)
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

	reader, writer := io.Pipe()
	g, ctx := errgroup.WithContext(r.Context())

	BesParts, goroutine := sa.FileInfo.FileManaging.FindBestOptions(sizeAndName.Size)

	timeS := time.Now()

	defer func() {
		sa := time.Since(timeS)
		fmt.Println(sa)
	}()

	chanelForAesKey := make(chan memguard.LockedBuffer)
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			defer func(writer *io.PipeWriter) {
				err := writer.Close()
				if err != nil {
					slog.Error("can't close a file", "err", err)
					return
				}
			}(writer)
			err = sa.EncryptFile(file, writer, chanelForAesKey)
			if err != nil {
				err := writer.CloseWithError(err)
				if err != nil {
					slog.Error("Error closing a file during encryption ", "Error", err)
					return err
				}
				return err
			}
			return nil
		}

	})

	GottenAesKey := <-chanelForAesKey

	defer GottenAesKey.Destroy()

	shortNameFile := sa.Crypto.Generate.GenerateShortName()
	FileExtension := sa.FileInfo.FileManaging.FindFormatOfFile(sizeAndName.Filename)

	Public, err := x509.ParsePKCS1PrivateKey(sa.Keys.ControllerKey.GetKey())
	if err != nil {
		slog.Error("Error parsing a key", "Error", err)
		err := writer.CloseWithError(err)
		if err != nil {
			slog.Error("Error in file writing wile a key parsing ", "Error", err)
			return "", err
		}
		err = reader.CloseWithError(err)
		if err != nil {
			slog.Error("Error in file reading wile a key parsing ", "Error", err)
			return "", err
		}

		return "", err
	}

	FileInfoInBytes, err := sa.Convert.Converting.JsonConverter(Dto.FileLabelsBytes{
		FileName: sizeAndName.Filename,
		AesKey:   hex.EncodeToString(GottenAesKey.Bytes()),
	})
	if err != nil {
		err := writer.CloseWithError(err)
		slog.Error("Error in file writing 2 ", "Error", err)
		return "", err
	}
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		errS3 := sa.S3.Uploader.UploadFileEncrypt(BesParts, goroutine, r.Context(), shortNameFile, FileExtension, reader)
		if errS3 != nil {
			//return "", err3
			return errS3
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	EncryptFileInfo, err := sa.Crypto.Encrypt.EncryptFileInfo(FileInfoInBytes, Public.Public().(*rsa.PublicKey))
	if err != nil {
		return "", err
	}

	err = sa.RedisControlling.Writer.WriteData(shortNameFile, EncryptFileInfo, r.Context())
	if err != nil {
		err := writer.CloseWithError(err)
		return "", err
	}

	time.AfterFunc(5*time.Minute, func() {
		g2, Ctx := errgroup.WithContext(context.Background())

		Ctx, cancel := context.WithTimeout(Ctx, 25*time.Second)
		defer cancel()
		DownloadingHaveStarted := sa.RedisControlling.CheckerRedis.ChekIsStartDownload(shortNameFile, Ctx)
		if DownloadingHaveStarted {
			return
		}
		g2.Go(func() error {

			err := sa.RedisControlling.Deleter.DeleteFileInfo(shortNameFile, Ctx)
			if err != nil {
				return err
			}
			return nil
		})
		g2.Go(func() error {
			err := sa.S3.Deleter.DeleteFileFromS3(shortNameFile, Ctx)
			if err != nil {
				return err
			}
			return nil
		})
		if err := g2.Wait(); err != nil {
			slog.Error("Error in file writing 2 ", "Error", err)
			return
		}
		slog.Info("Func Auto-deleteFile ends")
		return
	})

	slog.Info("File success upload ")

	return shortNameFile, nil

}

func (sa *HandlerPackCollect) FileUploadEncryptTest(fileName string) (string, error) {
	return "Hello World!", nil
}

func (sa *HandlerPackCollect) EncryptFile(file multipart.File, writer io.Writer, channelForBytes chan memguard.LockedBuffer) error {
	aesKey, err := memguard.NewBufferFromReader(rand.Reader, 32)
	if err != nil {
		slog.Error("Error generating random bytes", "Error", err)
		return err
	}
	defer aesKey.Destroy()

	go func() {
		channelForBytes <- *aesKey
		return
	}()

	block, err := aes.NewCipher(aesKey.Bytes())
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
			slog.Error("Error in file upload", "error", err.Error())
			return err
		}
		if err == io.EOF {
			break
		}
		stream.XORKeyStream(buf[:n], buf[:n])
		_, err = writer.Write(buf[:n])
		if err != nil {
			slog.Error("Err write in process", "Error", err.Error())
			return err
		}

	}

	return nil
}
