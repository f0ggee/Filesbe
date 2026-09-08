package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/FileControls"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"Kaban/internal/InfrastructureLayer/s3Repo"
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/awnumar/memguard"
	"golang.org/x/sync/errgroup"
)

type NewDownloadEncryptFileControl struct {
	Transfer FileControls.Transferring
}
type NewDownloadEncryptDelivery struct {
	ReaderRedis  DomainLevel.ReadingRedis
	DownloadS3   s3Repo.DownloadingS3
	DeleterRedis DomainLevel.DeleterRedis
	DeleterS3    s3Repo.DeleterS3
}

type NewDownloadEncryptCrypto struct {
	Decrypt       DomainLevel.Decrypter
	EncrypterKeys RepoEncrypterKeys.Keys
}
type NewDownloadEncrypt struct {
	NewDownloadEncryptDelivery
	NewDownloadEncryptCrypto
	NewDownloadEncryptFileControl
	d RepoParsers.Decode
}

func GetNewDownloadEncrypt(newDownloadEncryptDelivery NewDownloadEncryptDelivery, newDownloadEncryptCrypto NewDownloadEncryptCrypto, newDownloadEncryptFileControl NewDownloadEncryptFileControl) *NewDownloadEncrypt {
	return &NewDownloadEncrypt{NewDownloadEncryptDelivery: newDownloadEncryptDelivery, NewDownloadEncryptCrypto: newDownloadEncryptCrypto, NewDownloadEncryptFileControl: newDownloadEncryptFileControl}
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

	fileInfoInBytes, err := sa.ReaderRedis.GetFileInfo(data.EncryptedURl, data.Ctx)
	if err != nil {
		return err
	}
	safeFileData, err := sa.getFileData(fileInfoInBytes)
	if err != nil {
		return err
	}
	aesKey, fileName, err := sa.getFileInfoData(safeFileData.Data())
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
	Body, err := sa.DownloadS3.GetDownloadSecure(data.Ctx, data.EncryptedURl)
	if err != nil {
		return err
	}
	defer Body.Body.Close()
	g.Go(func() error {
		err = DecryptFile(aesKey.Data(), Body.Body, writer, ctx)
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
			err = sa.Transfer.TransferEncryptToClient(FileControls.TransferIncomingData{
				W: data.W,
				FileDetails: FileControls.FileDetails{
					FileFormat:   DomainLevel.GetNewFileSettings(0, fileName).FindFormatOfFile(),
					TrueFileName: fileName,
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
	err = sa.DeleterRedis.DeleteFileInfo(data.EncryptedURl, data.Ctx)
	if err != nil {
		return err
	}
	err = sa.DeleterS3.DeleteFileFromS3(data.EncryptedURl, data.Ctx)
	if err != nil {
		return err
	}
	return nil
}

func (sa *NewDownloadEncrypt) getFileInfoData(safeFileData []byte) (*memguard.LockedBuffer, string, error) {
	outputData := &DomainLevel.FileLabelsBytes{
		FileName: "",
		AesKey:   "",
	}
	err := sa.d.JsonDecodeMarshall(&sa, safeFileData)
	if err != nil {
		return nil, nil, errors.New(DomainLevel.ErrorParseInfo)
	}

	key := memguard.NewBuffer(len(outputData.AesKey))
	x, _ := hex.DecodeString(outputData.AesKey)
	key.Copy(x)
	return key, outputData.FileName, nil
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
			slog.Error("DecryptFile; error to read a File", "ERROR", err.Error())
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
func (sa *NewDownloadEncrypt) getFileData(fileInfoInBytes []byte) (*memguard.LockedBuffer, error) {
	outData, err := sa.Decrypt.DecryptData(DomainLevel.IncomeData{
		Key:  sa.EncrypterKeys.GetKey(),
		Data: fileInfoInBytes,
	})
	if err == nil {
		data := memguard.NewBuffer(len(outData))
		data.Copy(outData)
		memguard.WipeBytes(outData)
		return data, err
	}
	outData2, err := sa.Decrypt.DecryptData(DomainLevel.IncomeData{
		Key:  sa.EncrypterKeys.GetOldKey(),
		Data: fileInfoInBytes,
	})
	if err != nil {
		return nil, err
	}
	data := memguard.NewBuffer(len(outData2))
	data.Copy(outData2)
	memguard.WipeBytes(outData2)
	return data, nil
}
