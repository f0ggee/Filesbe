package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadEncryptRepo"
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type NewDownloadEncrypt struct {
	RedisControlling
	Crypto
	S3Controlling
	HandlerFileManagerPack
	FileDownloadEncrypt
	EncrypterKeys
}

func GetNewNewDownloadEncrypt(redisControlling RedisControlling, handlerPackCrypto Crypto, s3Controlling S3Controlling, handlerFileManagerPack HandlerFileManagerPack, fileDownloadEncrypt FileDownloadEncrypt, encrypterKeys EncrypterKeys) *NewDownloadEncrypt {
	return &NewDownloadEncrypt{RedisControlling: redisControlling, Crypto: handlerPackCrypto, S3Controlling: s3Controlling, HandlerFileManagerPack: handlerFileManagerPack, FileDownloadEncrypt: fileDownloadEncrypt, EncrypterKeys: encrypterKeys}
}

type NewDownloadEncryptNetwork struct {
	W http.ResponseWriter
}
type NewDownloadEncryptIncomingData struct {
	NewDownloadEncryptNetwork
	Ctx          context.Context
	EncryptedURl string
}

func (sa *NewDownloadEncrypt) DownloadEncrypt(data NewDownloadEncryptIncomingData) error {

	fileInfoInBytes, err := sa.RedisControlling.Reader.GetFileInfo(data.EncryptedURl, data.Ctx)
	if err != nil {
		return err
	}

	aesKey, realFileName, err := sa.Decrypt.DecryptFileInfo(fileInfoInBytes, sa.GetKeys.GetKey(), sa.GetKeys.GetOldKey())
	if err != nil {
		return err
	}

	Reader, writer := io.Pipe()
	defer func(writer *io.PipeWriter) {
		err := writer.Close()
		if err != nil {
			slog.Error("Writer can't close", "Err", err)
			return
		}
	}(writer)
	g, ctx := errgroup.WithContext(data.Ctx)
	Body, err := sa.S3Download.GetDownloadSecure(data.Ctx, data.EncryptedURl)
	if err != nil {
		return err
	}
	defer Body.Body.Close()
	g.Go(func() error {
		err = DecryptFile(aesKey, Body.Body, writer, ctx)
		if err != nil {
			err := writer.CloseWithError(err)
			if err != nil {
				slog.Error("DownloadEncrypt; error to close the body", "ERROR", err)
				return err
			}
		}
		return nil
	})

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return data.Ctx.Err()
		default:
			err = sa.FileDownloadEncrypt.EncryptDownload.EncryptDownloadFile(RepoDownloadEncryptRepo.DownloadEncryptIncomingData{
				W: data.W,
				FileDetails: RepoDownloadEncryptRepo.FileDetails{
					FileFormat:   sa.FileManaging.FindFormatOfFile(realFileName),
					TrueFileName: realFileName,
					FileLength:   *Body.ContentLength - aes.BlockSize,
				},
				FileBody: Reader,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
	err = sa.RedisControlling.Deleter.DeleteFileInfo(data.EncryptedURl, data.Ctx)
	if err != nil {
		return err
	}
	err = sa.S3Controlling.Deleter.DeleteFileFromS3(data.EncryptedURl, data.Ctx)
	if err != nil {
		return err
	}
	return nil
}
func DecryptFile(AesKey []byte, o io.ReadCloser, writer *io.PipeWriter, ctx context.Context) error {
	block, err := aes.NewCipher(AesKey)
	if err != nil {
		slog.Error("DecryptFile; error to make a new AES decrypting", "ERROR", err.Error())
		return errors.New(DomainLevel.ErrorDecryptFile)
	}
	nonce := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(o, nonce)
	if err != nil {
		slog.Error("DecryptFile; there is error to generate a nonce", "ERROR", err.Error())
		return errors.New(DomainLevel.ErrorDecryptFile)
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
			slog.Error("DecryptFile; error to read a file", "ERROR", err.Error())
			return errors.New(DomainLevel.ErrorDecryptFile)
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
					slog.Error("DecryptFile: error to write into a stream", "ERROR", err)
					return errors.New(DomainLevel.ErrorDecryptFile)
				}
				return errors.New(DomainLevel.ErrorDecryptFile)
			}
		}
	}
	return nil
}
