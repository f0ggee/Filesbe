package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/FileControls"
	s3Repo2 "Kaban/internal/InfrastructureLayer/FileTransferring/s3Repo"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"bufio"
	"context"
	"crypto/aes"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/awnumar/memguard"
	"golang.org/x/sync/errgroup"
)

const ErrorDecryptFile = "error to decrypt data"

type NewDownloadEncryptFileControl struct {
	Transfer FileControls.Transferring
}
type NewDownloadEncryptDelivery struct {
	ReaderRedis  DomainLevel.ReadingRedis
	DownloadS3   s3Repo2.DownloadingS3
	DeleterRedis DomainLevel.DeleterRedis
	DeleterS3    s3Repo2.DeleterS3
}

type NewDownloadEncryptCrypto struct {
	Decrypt       DomainLevel.CryptoMaker
	DecryptSpec   DomainLevel.CryptoMakerSpec
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

func (s *NewDownloadEncrypt) DownloadEncrypt(data NewDownloadEncryptIncomingData) error {
	Reader, writer := io.Pipe()
	Body, err := s.DownloadS3.GetDownloadSecure(data.Ctx, data.EncryptedURl)
	defer func() {
		err = s.closeSources(Body.Body, writer, Reader)
		if err != nil {
			return
		}
	}()
	fileInfoInBytes, err := s.ReaderRedis.GetFileInfo(data.EncryptedURl, data.Ctx)
	if err != nil {
		return err
	}
	safeFileData, err := s.getFileData(fileInfoInBytes)
	if err != nil {
		return err
	}
	aesKey, fileName, err := s.getFileInfoData(safeFileData.Data())
	if err != nil {
		return err
	}
	defer aesKey.Destroy()

	g, ctx := errgroup.WithContext(data.Ctx)
	if err != nil {
		return err
	}
	s.setStartDecrypter(g, aesKey.Bytes(), Body.Body, writer, ctx)
	s.setGoroutineDownloader(data, g, ctx, FileControls.TransferIncomingData{
		W: data.W,
		FileDetails: FileControls.FileDetails{
			FileFormat:   DomainLevel.GetNewFileSettings(0, fileName).FindFormatOfFile(),
			TrueFileName: fileName,
			FileLength:   *Body.ContentLength - aes.BlockSize,
		},
		FileBody: Reader,
	})
	if err = g.Wait(); err != nil {
		return err
	}
	err = s.DeleterRedis.DeleteFileInfo(data.EncryptedURl, data.Ctx)
	if err != nil {
		return err
	}
	err = s.DeleterS3.DeleteFileFromS3(data.EncryptedURl, data.Ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *NewDownloadEncrypt) setGoroutineDownloader(data NewDownloadEncryptIncomingData, g *errgroup.Group, ctx context.Context, d FileControls.TransferIncomingData) {
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return data.Ctx.Err()
		default:
			err := s.Transfer.TransferEncryptToClient(d)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *NewDownloadEncrypt) setStartDecrypter(g *errgroup.Group, aesKey []byte, Body io.Reader, writer *io.PipeWriter, ctx context.Context) {
	g.Go(func() error {
		err := s.decryptFile(aesKey, Body, writer, ctx)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *NewDownloadEncrypt) getFileInfoData(safeFileData []byte) (*memguard.LockedBuffer, string, error) {
	var outputData DomainLevel.FileLabelsBytes
	err := s.d.JsonDecodeMarshall(&s, safeFileData)
	if err != nil {
		return nil, "", err
	}
	key := memguard.NewBuffer(len(outputData.AesKey))
	x, err := hex.DecodeString(outputData.AesKey)
	if err != nil {
		return nil, "", errors.New(DomainLevel.ErrorParseInfo)
	}
	key.Copy(x)
	memguard.WipeBytes(x)
	return key, outputData.FileName, nil
}

func (s *NewDownloadEncrypt) decryptFile(AesKey []byte, o io.Reader, writer *io.PipeWriter, ctx context.Context) error {
	nonce := make([]byte, aes.BlockSize+len(AesKey))
	_, err := io.ReadFull(o, nonce[:aes.BlockSize])
	if err != nil {
		slog.Error("decryptFile; there is error to generate a nonce", "ERROR", err.Error())
		return errors.New(ErrorStartUploading)
	}
	nonce = append(nonce, AesKey...)
	decrBlock, err := s.Decrypt.MakeCrypto(nonce)
	if err != nil {
		return err
	}
	plainText := memguard.NewBuffer(35 * 1024)
	defer plainText.Destroy()
	file := bufio.NewReader(o)
	for {
		if ctx.Err() != nil {
			return errors.New("context canceled")
		}
		n, err := file.Read(plainText.Bytes())
		if err != nil {
			slog.Error("decryptFile; error to read a File", "ERROR", err.Error())
			return errors.New(ErrorDecryptFile)
		}
		if err == io.EOF {

			break
		}
		if n > 0 {
			decrData, err := decrBlock.Decrypt(plainText.Bytes()[:n])
			_, err = writer.Write(decrData)
			if err != nil {
				slog.Error("decryptFile: error to write into a stream", "ERROR", err)
				return errors.New(ErrorDecryptFile)
			}
		}
	}
	return nil
}
func (s *NewDownloadEncrypt) getFileData(fileInfoInBytes []byte) (*memguard.LockedBuffer, error) {
	cr, err := s.DecryptSpec.MakeCrypto(s.EncrypterKeys.GetKey(), 0)
	if err != nil {
		return nil, err
	}
	OutData, err := cr.Decrypt(fileInfoInBytes)
	if err == nil {
		data := memguard.NewBuffer(len(OutData))
		data.Copy(OutData)
		memguard.WipeBytes(OutData)
		return data, nil
	}
	cr, err = s.DecryptSpec.MakeCrypto(s.EncrypterKeys.GetOldKey(), 0)
	if err != nil {
		return nil, errors.New(ErrorDecryptFile)
	}
	outData2, err := cr.Decrypt(fileInfoInBytes)
	if err != nil {
		return nil, errors.New(ErrorDecryptFile)
	}
	data := memguard.NewBuffer(len(outData2))
	data.Copy(outData2)
	memguard.WipeBytes(outData2)
	return data, nil
}
func (s *NewDownloadEncrypt) closeSources(reader io.ReadCloser, writer *io.PipeWriter, reader2 *io.PipeReader) error {

	err := reader.Close()
	if err != nil {
		slog.Error("DownloadEncrypt; error to close a writer", "ERROR", err)
		return err
	}
	err = writer.Close()
	if err != nil {
		slog.Error("DownloadEncrypt; error to close a writer", "ERROR", err)
		return err
	}
	err = reader2.Close()
	if err != nil {
		slog.Error("DownloadEncrypt; error to close a reader", "ERROR", err)
		return err
	}
	return nil
}
