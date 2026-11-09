package Application

import (
	"Kaban/internal/DomainLevel"
	"bytes"
	"context"
	"crypto/rsa"
	"errors"
	"io"
	"log/slog"

	"github.com/awnumar/memguard"
	"golang.org/x/sync/errgroup"
)

type NewDownloadEncryptFileControl struct {
	Uploader   DomainLevel.MakerUploader
	Downloader DomainLevel.MakerDownloader
	Deleter    DomainLevel.MakerDeleter
}
type NewDownloadEncryptDelivery struct {
	ReaderRedis  DomainLevel.ReadingRedis
	DeleterRedis DomainLevel.DeleterRedis
}

type NewDownloadEncryptCrypto struct {
	Aes DomainLevel.CryptoMaker
	Rsa DomainLevel.CryptoMaker
}
type DownloadEncrypt struct {
	NewDownloadEncryptDelivery
	NewDownloadEncryptCrypto
	NewDownloadEncryptFileControl
	Decoder DomainLevel.Decoder
}

type DownloadEncryptIncomeData struct {
	dst           io.Writer
	Ctx           context.Context
	EncryptedName string
}

func (s *DownloadEncrypt) DownloadEncryptApplication(data DownloadEncryptIncomeData) error {

	g, _ := errgroup.WithContext(data.Ctx)

	fileInfoInBytes, err := s.ReaderRedis.GetFileInfo(data.EncryptedName, data.Ctx)
	if err != nil {
		return err
	}

	key, fileName, err := s.getFileData(fileInfoInBytes)
	if err != nil {
		return err
	}
	defer func() {
		key.Destroy()
		fileName.Destroy()
	}()
	downloadedObject, err := s.Downloader.SetName(data.EncryptedName).Make(data.Ctx)
	if err != nil {
		return err
	}
	defer downloadedObject.CloseSource()

	g.Go(func() error {
		uploadedObject, err := s.Uploader.SetName(fileName.String()).SetAdditionalWriter(data.dst).SetSize(s.Downloader.GetFileSize()).Make(data.Ctx)
		if err != nil {
			return err
		}
		defer uploadedObject.CloseSource()
		err = s.decryptFile(key.Bytes(), downloadedObject, uploadedObject)
		if err != nil {
			return err
		}
		return nil
	})
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

func (s *DownloadEncrypt) decryptFile(AesKey []byte, src DomainLevel.Download, dst DomainLevel.Upload) error {
	download, err := src.Downloader()

	if err != nil {
		return err
	}

	nonce, err := memguard.NewBufferFromReader(download, s.Aes.GetRequiredOverheadSize())
	if err != nil {
		slog.Error("decryptFile; there is error to generate a nonce", "ERROR", err.Error())
		return ErrorDecryptFile
	}
	decrBlock, err := s.Aes.MakeCrypto(AesKey[s.Aes.GetRequiredOverheadSize():], nonce.Bytes())
	if err != nil {
		return err
	}
	plainText := memguard.NewBuffer(ChunkSize)
	defer plainText.Destroy()
	for {
		n, err := download.Read(plainText.Bytes())
		if n > 0 {
			decrData, err := decrBlock.Decrypt(plainText.Bytes()[:n])
			if err != nil {
				slog.Error("decryptFile: error to write into a stream", "ERROR", err)
				return ErrorDecryptFile
			}
			err = dst.Upload(bytes.NewReader(decrData))
			if err != nil {
				slog.Error("decryptFile: error to write into a stream", "ERROR", err)
				memguard.WipeBytes(decrData)
				return ErrorDecryptFile
			}
			memguard.WipeBytes(decrData)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("decryptFile; error to read a File", "ERROR", err.Error())
			return ErrorDecryptFile
		}

	}
	return nil
}

func (s *DownloadEncrypt) getFileData(fileInfoInBytes []byte) (*memguard.LockedBuffer, *memguard.LockedBuffer, error) {
	outputData := &memguard.LockedBuffer{}
	defer outputData.Destroy()

	cr, errMakeCrypto := s.Rsa.MakeCrypto(DomainLevel.GetCryptoKeys.GetKey(), []byte(""))
	if errMakeCrypto != nil {
		return nil, nil, ErrorDecryptFile
	}
	outDataNewKey, err := cr.Decrypt(fileInfoInBytes)
	switch {
	case errors.Is(err, rsa.ErrDecryption):
		cr, err = s.Rsa.MakeCrypto(DomainLevel.GetCryptoKeys.GetOldKey(), []byte(""))
		if err != nil {
			return nil, nil, ErrorDecryptFile
		}
		outDataOldKey, err := cr.Decrypt(fileInfoInBytes)
		if err != nil {
			return nil, nil, ErrorDecryptFile
		}
		outputData = memguard.NewBufferFromBytes(outDataOldKey)
	default:
		outputData = memguard.NewBufferFromBytes(outDataNewKey)
	}

	var fileData DomainLevel.FileLabelsBytes
	err = s.Decoder.Decode(&fileData, outputData.Bytes())
	if err != nil {
		return nil, nil, err
	}
	key := memguard.NewBufferFromBytes(fileData.AesKey)
	fileName := memguard.NewBufferFromBytes(fileData.FileName)
	return key, fileName, nil
}
