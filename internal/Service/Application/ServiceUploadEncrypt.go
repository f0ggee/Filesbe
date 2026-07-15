package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/FileControls"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
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

type NewUploadEncryptCrypto struct {
	Generate   DomainLevel.CryptoGenerating
	ServerKeys DomainLevel.NewServerKeys
	Encrypt    DomainLevel.Encryption
}
type NewUploadEncryptDataManage struct {
	FileManger FileControls.FileSettings
	Encode     RepoParsers.Encode
}
type NewUploadEncryptDelivery struct {
	UploaderS3   s3Repo.S3Uploader
	RedisWriter  DomainLevel.WritingRedis
	RedisChecker DomainLevel.RedisChecker
	RedisDeleter DomainLevel.DeleterRedis
	DeleterS3    s3Repo.DeleterS3
}
type NewUploadEncrypt struct {
	NewUploadEncryptDataManage
	NewUploadEncryptCrypto
	NewUploadEncryptDelivery
}

func GetNewUploadEncrypt(newUploadEncryptDataManage NewUploadEncryptDataManage, newUploadEncryptCrypto NewUploadEncryptCrypto, newUploadEncryptDelivery NewUploadEncryptDelivery) *NewUploadEncrypt {
	return &NewUploadEncrypt{NewUploadEncryptDataManage: newUploadEncryptDataManage, NewUploadEncryptCrypto: newUploadEncryptCrypto, NewUploadEncryptDelivery: newUploadEncryptDelivery}
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
			slog.Error("UploadEncrypt; error to close a file flow", "ERROR", err)
			return
		}
	}()

	reader, writer := io.Pipe()
	g, ctx := errgroup.WithContext(r.Context())

	BesParts, goroutine := sa.FileManger.FindBestOptions(sizeAndName.Size)

	chanelForAesKey := make(chan memguard.LockedBuffer)
	sa.startEncryptFile(g, ctx, writer, file, chanelForAesKey)
	GottenAesKey := <-chanelForAesKey

	defer GottenAesKey.Destroy()

	shortNameFile := sa.Generate.GenerateShortName()
	FileExtension := sa.FileManger.FindFormatOfFile(sizeAndName.Filename)

	OurKey := sa.ServerKeys.GerOurPrivateKey()
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

	FileInfoInBytes, err := sa.Encode.JsonEncodeMarshall(DomainLevel.FileLabelsBytes{
		FileName: sizeAndName.Filename,
		AesKey:   hex.EncodeToString(GottenAesKey.Bytes()),
	})
	if err != nil {
		slog.Error("UploadEncrypt; error to encode data into Json", "ERROR", err)
		writer.CloseWithError(err)
		return "", errors.New(DomainLevel.ErrorStrangeUploadFile)
	}

	sa.setFileUploader(g, ctx, BesParts, goroutine, FileExtension, shortNameFile, reader)
	if err := g.Wait(); err != nil {
		return "", err
	}

	EncryptFileInfo, err := sa.Encrypt.EncryptFileInfo(FileInfoInBytes, Public.Public().(*rsa.PublicKey))
	if err != nil {
		return "", err
	}

	err = sa.RedisWriter.WriteData(DomainLevel.WriteDataIncomeData{
		FileName: shortNameFile,
		Info:     EncryptFileInfo,
		Ctx:      r.Context(),
	})
	if err != nil {
		err := writer.CloseWithError(err)
		return "", err
	}

	sa.setDeleteFile(shortNameFile)

	return shortNameFile, nil

}

func (sa *NewUploadEncrypt) setFileUploader(g *errgroup.Group, ctx context.Context, BesParts int, goroutine int, FileExtension string, shortNameFile string, reader *io.PipeReader) {
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		errS3 := sa.UploaderS3.UploadFileEncrypt(s3Repo.UploadFileIncomingData{
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
			return errS3
		}
		return nil
	})
}

func (sa *NewUploadEncrypt) setDeleteFile(shortNameFile string) *time.Timer {
	return time.AfterFunc(5*time.Minute, func() {
		g2, Ctx := errgroup.WithContext(context.Background())
		Ctx, cancel := context.WithTimeout(Ctx, 25*time.Second)
		defer cancel()
		DownloadingHaveStarted := sa.RedisChecker.ChekIsStartDownload(shortNameFile, Ctx)
		if DownloadingHaveStarted {
			return
		}
		g2.Go(func() error {

			err := sa.RedisDeleter.DeleteFileInfo(shortNameFile, Ctx)
			if err != nil {
				return err
			}
			return nil
		})
		g2.Go(func() error {
			err := sa.DeleterS3.DeleteFileFromS3(shortNameFile, Ctx)
			if err != nil {
				return err
			}
			return nil
		})
		if err := g2.Wait(); err != nil {
			slog.Error("DeleterS3; error to delete a file", "ERROR", err)
			return
		}
		return
	})
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
	buf := memguard.NewBuffer(32 * 1024)
	defer buf.Destroy()
	_, err = writer.Write(nonce)
	if err != nil {
		slog.Error("EncryptFile; error happened during writing", "ERROR", err)
		return errors.New(DomainLevel.ErrorStartUploading)
	}
	for {
		n, err := file.Read(buf.Bytes())
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("EncryptFile; error to download a file", "ERROR", err.Error())
			return errors.New(DomainLevel.ErrorStrangeUploadFile)
		}
		stream.XORKeyStream(buf.Bytes()[:n], buf.Bytes()[:n])
		_, err = writer.Write(buf.Bytes()[:n])
		if err != nil {
			slog.Error("EncryptFile; error to write in the stream ", "ERROR", err.Error())
			return errors.New(DomainLevel.ErrorStrangeCrypto)
		}
	}

	return nil
}
