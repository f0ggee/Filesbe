package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/FileControls"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"bufio"
	"bytes"
	"context"
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
	Transfer   FileControls.Transferring
	Uploader   DomainLevel.MakerUploader
	Downloader DomainLevel.MakerDownloader
	Deleter    DomainLevel.MakerDeleter
}
type NewDownloadEncryptDelivery struct {
	ReaderRedis  DomainLevel.ReadingRedis
	DeleterRedis DomainLevel.DeleterRedis
}

type NewDownloadEncryptCrypto struct {
	Crypto        DomainLevel.CryptoMaker
	EncrypterKeys RepoEncrypterKeys.Keys
}
type NewDownloadEncrypt struct {
	NewDownloadEncryptDelivery
	NewDownloadEncryptCrypto
	NewDownloadEncryptFileControl
	d DomainLevel.Decoder
}

func GetNewDownloadEncrypt(newDownloadEncryptDelivery NewDownloadEncryptDelivery, newDownloadEncryptCrypto NewDownloadEncryptCrypto, newDownloadEncryptFileControl NewDownloadEncryptFileControl) *NewDownloadEncrypt {
	return &NewDownloadEncrypt{NewDownloadEncryptDelivery: newDownloadEncryptDelivery, NewDownloadEncryptCrypto: newDownloadEncryptCrypto, NewDownloadEncryptFileControl: newDownloadEncryptFileControl}
}

type NewDownloadEncryptNetwork struct {
	W http.ResponseWriter
}
type NewDownloadEncryptIncomingData struct {
	NewDownloadEncryptNetwork
	Ctx           context.Context
	EncryptedName string
}

func (s *NewDownloadEncrypt) DownloadEncrypt(data NewDownloadEncryptIncomingData) error {
	Reader, writer := io.Pipe()
	g, ctx := errgroup.WithContext(data.Ctx)
	downloadedObject, err := s.Downloader.SetName(data.EncryptedName).Make(data.Ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = s.closeSources(nil, writer, Reader)
		if err != nil {
			return
		}
	}()

	readyDownload, err := downloadedObject.Downloader()
	if err != nil {
		return err
	}
	fileInfoInBytes, err := s.ReaderRedis.GetFileInfo(data.EncryptedName, data.Ctx)
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
	uploadedObject, err := s.Uploader.SetName(fileName).SetAdditionalWriter(data.W).SetSize(s.Downloader.GetFileSize()).Make(data.Ctx)
	if err != nil {
		return err
	}

	g.Go(func() error {
		err = s.decryptFile(ctx, aesKey.Bytes(), readyDownload, uploadedObject)
		if err != nil {
			return err
		}
		return nil
	})
	//s.setStartDecrypter(g, aesKey.Bytes(), readyDownload, writer, ctx)
	s.setGoroutineDownloader(g, ctx, uploadedObject, Reader)
	if err = g.Wait(); err != nil {
		return err
	}
	err = s.DeleterRedis.DeleteFileInfo(data.EncryptedName, data.Ctx)
	if err != nil {
		return err
	}
	deleteObject, err := s.Deleter.SetName(data.EncryptedName).Make(data.Ctx)
	if err != nil {
		return err
	}
	return deleteObject.Deleter()
}

func (s *NewDownloadEncrypt) setGoroutineDownloader(g *errgroup.Group, ctx context.Context, d DomainLevel.Upload, read io.Reader) {
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := d.Uploader(read)
			if err != nil {
				return err
			}
			return nil
		}
	})
}

func (s *NewDownloadEncrypt) getFileInfoData(safeFileData []byte) (*memguard.LockedBuffer, string, error) {
	var outputData DomainLevel.FileLabelsBytes
	err := s.d.Decode(&s, safeFileData)
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

func (s *NewDownloadEncrypt) decryptFile(ctx context.Context, AesKey []byte, o io.Reader, uploadedObj DomainLevel.Upload) error {

	nonce := memguard.NewBuffer(s.Crypto.GetRequiredRandomSize() + len(AesKey))
	_, err := io.ReadFull(o, nonce.Bytes()[:s.Crypto.GetRequiredRandomSize()])
	if err != nil {
		slog.Error("decryptFile; there is error to generate a nonce", "ERROR", err.Error())
		return errors.New(ErrorStartUploading)
	}
	copy(nonce.Bytes()[s.Crypto.GetRequiredRandomSize():], AesKey)
	decrBlock, err := s.Crypto.MakeCrypto(nonce.Bytes()[s.Crypto.GetRequiredRandomSize():], nonce.Bytes()[:s.Crypto.GetRequiredRandomSize()])
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
		if n > 0 {
			decrData, err := decrBlock.Decrypt(plainText.Bytes()[:n])
			if err != nil {
				slog.Error("decryptFile: error to write into a stream", "ERROR", err)
				return errors.New(ErrorDecryptFile)
			}
			err = uploadedObj.Uploader(bytes.NewReader(decrData))
			if err != nil {
				slog.Error("decryptFile: error to write into a stream", "ERROR", err)
				return errors.New(ErrorDecryptFile)
			}
		}
		if err == io.EOF {

			break
		}
		if err != nil {
			slog.Error("decryptFile; error to read a File", "ERROR", err.Error())
			return errors.New(ErrorDecryptFile)
		}

	}
	return nil
}
func (s *NewDownloadEncrypt) getFileData(fileInfoInBytes []byte) (*memguard.LockedBuffer, error) {
	cr, err := s.Crypto.MakeCrypto(s.EncrypterKeys.GetKey(), []byte(""))
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
	cr, err = s.Crypto.MakeCrypto(s.EncrypterKeys.GetOldKey(), []byte(""))
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
	return nil, nil
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
