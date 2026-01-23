package Handlers

import (
	Dto2 "Kaban/iternal/Dto"
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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func DownloadEncrypt(w http.ResponseWriter, ctxs context.Context, name string) error {

	ctx, cancel := Uttiltesss2.ContextForDownloading(ctxs)
	defer cancel()

	//It,saves the file name in the variable for comfortably
	NameOfFile := name

	NameOfFile = GetTrueName(NameOfFile, name)

	sa, err2 := PrecessingFile(NameOfFile)
	if err2 != nil {

		return errors.New("file was used")
	}

	//Return the string as bytes
	ReturnToBytes, err := hex.DecodeString(sa.AesKey)
	if err != nil {
		slog.Error("Error decode to string", "Err", err)
		return err
	}

	Mut.RLock()
	PrivatKey := NewPrivateKey
	Mut.RUnlock()
	//Decrypt the key
	AesDecryptKey, err := rsa.DecryptOAEP(sha256.New(), nil, PrivatKey, ReturnToBytes, nil)

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
		Key:         aws.String(NameOfFile),
	})

	defer o.Body.Close()

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
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename= %v", NameOfFile))
	w.Header().Set("Content-Length", strconv.FormatInt(*o.ContentLength-aes.BlockSize, 10))

	if _, err = io.Copy(w, Reader); err != nil {
		slog.Error("Err In file Service Downloader Encrypt", "err", err)
		return errors.New("connect close")
	}

	slog.Info("Func Download ends")
	//Here, we start the func, which deletes a file
	DeleteFile(NameOfFile, false)

	return nil
}

func PrecessingFile(NameOfFile string) (struct {
	AesKey          string
	TimeSet         time.Time
	IsUsed          bool
	IsStartDownload bool
}, error) {

	Mut.RLock()
	sa, ok := Dto2.MapForFile[NameOfFile]
	Mut.RUnlock()

	if !ok {
		slog.Error("File don't find ")
		return struct {
			AesKey          string
			TimeSet         time.Time
			IsUsed          bool
			IsStartDownload bool
		}{}, errors.New("don't find ")
	}

	//If a file is used
	if sa.IsUsed {
		return struct {
			AesKey          string
			TimeSet         time.Time
			IsUsed          bool
			IsStartDownload bool
		}{}, errors.New("file was used")
	}

	//Here, we mention that a file starts downloads
	sa.IsStartDownload = true
	return sa, nil
}

func GetTrueName(NameOfFile string, name string) string {
	Mut.RLock()
	if Name, ok := Dto2.NamesToConvert[NameOfFile]; ok {
		NameOfFile = Name
		delete(Dto2.NamesToConvert, name)
	}
	Mut.RUnlock()
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
