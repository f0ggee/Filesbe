package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/s3Repo"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/awnumar/memguard"
	"golang.org/x/sync/errgroup"
)

type NewUploadEncrypt struct {
	GetFileManager
	GetCrypto
	Parser
	GetControlKeys
	S3Controlling
	RedisControlling
}

func GetNewNewUploadEncrypt(getFileManager GetFileManager, getCrypto GetCrypto, parser Parser, getControlKeys GetControlKeys, s3Controlling S3Controlling, redisControlling RedisControlling) *NewUploadEncrypt {
	return &NewUploadEncrypt{GetFileManager: getFileManager, GetCrypto: getCrypto, Parser: parser, GetControlKeys: getControlKeys, S3Controlling: s3Controlling, RedisControlling: redisControlling}
}

func (sa *NewUploadEncrypt) UploadEncrypt(r *http.Request) (string, error) {

	file, sizeAndName, err := r.FormFile("file")
	if err != nil {
		slog.Error("UploadEncrypt; error to get a file", "ERROR", err)
		return "", errors.New(DomainLevel.ErrorStartUploading)
	}
	if sizeAndName.Size >= DomainLevel.FileMaxSize {
		return "", errors.New(DomainLevel.ErrorFileSizeBig)
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

	BesParts, goroutine := sa.FileManaging.FindBestOptions(sizeAndName.Size)

	chanelForAesKey := make(chan memguard.LockedBuffer)
	sa.startEncryptFile(g, ctx, writer, file, chanelForAesKey)
	GottenAesKey := <-chanelForAesKey

	defer GottenAesKey.Destroy()

	shortNameFile := sa.Generate.GenerateShortName()
	FileExtension := sa.FileManaging.FindFormatOfFile(sizeAndName.Filename)

	OurKey := sa.GetControlKeys.Keys.GerOurPrivateKey()
	Public, err := x509.ParsePKCS1PrivateKey(OurKey)
	if err != nil {
		slog.Error("UploadEncrypt; error to decode a key", "ERROR", err)
		err := writer.CloseWithError(err)
		if err != nil {
			slog.Error("Error in file writing wile a key parsing ", "Error", err)
			return "", errors.New(DomainLevel.ErrorStrangeUploadFile)
		}
		err = reader.CloseWithError(err)
		if err != nil {
			slog.Error("Error in file reading wile a key parsing ", "Error", err)
			return "", errors.New(DomainLevel.ErrorStrangeUploadFile)
		}
		return "", errors.New(DomainLevel.ErrorStrangeUploadFile)
	}

	FileInfoInBytes, err := sa.Parser.Encode.JsonEncodeMarshall(DomainLevel.FileLabelsBytes{
		FileName: sizeAndName.Filename,
		AesKey:   hex.EncodeToString(GottenAesKey.Bytes()),
	})
	if err != nil {
		slog.Error("UploadEncrypt; error to encode data into Json", "ERROR", err)
		writer.CloseWithError(err)
		return "", errors.New(DomainLevel.ErrorStrangeUploadFile)
	}

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		errS3 := sa.S3Controlling.Uploader.UploadFileEncrypt(s3Repo.UploadFileIncomingData{
			Parts:      BesParts,
			Goroutines: goroutine,
			Ctx:        ctx,
			FileDetails: s3Repo.FileDetails{
				FileFormat: FileExtension,
				FileName:   shortNameFile,
				FileBody: s3Repo.TypeUploading{
					Pipe: reader,
				},
			},
		})
		if errS3 != nil {
			//return "", err3
			return errS3
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return "", err
	}

	EncryptFileInfo, err := sa.Encrypt.EncryptFileInfo(FileInfoInBytes, Public.Public().(*rsa.PublicKey))
	if err != nil {
		return "", err
	}

	err = sa.Writer.WriteData(shortNameFile, EncryptFileInfo, r.Context())
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
			err := sa.S3Controlling.Deleter.DeleteFileFromS3(shortNameFile, Ctx)
			if err != nil {
				return err
			}
			return nil
		})
		if err := g2.Wait(); err != nil {
			slog.Error("Deleter; error to delete a file", "ERROR", err)
			return
		}
		return
	})

	return shortNameFile, nil

}

func (sa *NewUploadEncrypt) startEncryptFile(g *errgroup.Group, ctx context.Context, writer *io.PipeWriter, file multipart.File, chanelForAesKey chan memguard.LockedBuffer) {
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
			err := sa.EncryptFile(file, writer, chanelForAesKey)
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
}

func (sa *NewUploadEncrypt) EncryptFile(file multipart.File, writer io.Writer, channelForBytes chan memguard.LockedBuffer) error {
	aesKey, err := memguard.NewBufferFromReader(rand.Reader, 32)
	if err != nil {
		slog.Error("EncryptFile; error happened during generating the aes Key", "ERROR", err)
		return errors.New(DomainLevel.ErrorStrangeCrypto)
	}
	defer aesKey.Destroy()
	go func() {
		channelForBytes <- *aesKey
		return
	}()

	block, err := aes.NewCipher(aesKey.Bytes())
	if err != nil {
		slog.Error("EncryptFile; error to make a new cipher block", "ERROR", err)
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
		slog.Error("EncryptFile; error happened during writing", "ERROR", err)
		return err
	}
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("EncryptFile; error to download a file", "ERROR", err.Error())
			return errors.New(DomainLevel.ErrorStrangeUploadFile)
		}
		stream.XORKeyStream(buf[:n], buf[:n])
		_, err = writer.Write(buf[:n])
		if err != nil {
			slog.Error("EncryptFile; error to write in the stream ", "ERROR", err.Error())
			return errors.New(DomainLevel.ErrorStrangeCrypto)
		}
	}

	return nil
}
