package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"bytes"
	"context"
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
	Aes           DomainLevel.CryptoMaker
	Rsa           DomainLevel.CryptoMaker
	EncrypterKeys RepoEncrypterKeys.Keys
}
type DownloadEncrypt struct {
	NewDownloadEncryptDelivery
	NewDownloadEncryptCrypto
	NewDownloadEncryptFileControl
	decoder DomainLevel.Decoder
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

	safeFileData, err := s.getFileData(fileInfoInBytes)
	if err != nil {
		return err
	}
	defer safeFileData.Destroy()
	aesKey, fileName, err := s.getFileInfoData(safeFileData.Data())
	if err != nil {
		return err
	}
	defer aesKey.Destroy()

	downloadedObject, err := s.Downloader.SetName(fileName).Make(data.Ctx)
	if err != nil {
		return err
	}
	defer downloadedObject.CloseSource()
	defer aesKey.Destroy()

	g.Go(func() error {
		uploadedObject, err := s.Uploader.SetName(fileName).SetAdditionalWriter(data.dst).SetSize(s.Downloader.GetFileSize()).Make(data.Ctx)
		if err != nil {
			return err
		}
		defer uploadedObject.CloseSource()
		err = s.decryptFile(aesKey.Bytes(), downloadedObject, uploadedObject)
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

func (s *DownloadEncrypt) getFileInfoData(safeFileData []byte) (*memguard.LockedBuffer, string, error) {
	var outputData DomainLevel.FileLabelsBytes
	err := s.decoder.Decode(&outputData, safeFileData)
	if err != nil {
		return nil, "", err
	}
	key := memguard.NewBuffer(len(outputData.AesKey))
	key.Copy(outputData.AesKey)
	memguard.WipeBytes(outputData.AesKey)
	return key, outputData.FileName, nil
}

func (s *DownloadEncrypt) decryptFile(AesKey []byte, o DomainLevel.Download, uploadedObj DomainLevel.Upload) error {
	download, err := o.Downloader()
	if err != nil {
		return err
	}
	nonce := memguard.NewBuffer(s.Aes.GetRequiredOverheadSize())
	_, err = io.ReadFull(download, nonce.Bytes())
	if err != nil {
		slog.Error("decryptFile; there is error to generate a nonce", "ERROR", err.Error())
		return ErrorDecryptFile
	}
	decrBlock, err := s.Aes.MakeCrypto(AesKey, nonce.Bytes())
	if err != nil {
		return err
	}
	plainText := memguard.NewBuffer(DomainLevel.ChunkSize)
	defer plainText.Destroy()
	for {
		n, err := download.Read(plainText.Bytes())
		if n > 0 {
			decrData, err := decrBlock.Decrypt(plainText.Bytes()[:n])
			if err != nil {
				slog.Error("decryptFile: error to write into a stream", "ERROR", err)
				return ErrorDecryptFile
			}
			err = uploadedObj.Uploader(bytes.NewReader(decrData))
			if err != nil {
				slog.Error("decryptFile: error to write into a stream", "ERROR", err)
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
func (s *DownloadEncrypt) getFileData(fileInfoInBytes []byte) (*memguard.LockedBuffer, error) {
	var data *memguard.LockedBuffer
	cr, err := s.Rsa.MakeCrypto(s.EncrypterKeys.GetKey(), []byte(""))
	if err != nil {
		return nil, err
	}
	OutData, err := cr.Decrypt(fileInfoInBytes)
	if err == nil {
		data = memguard.NewBuffer(len(OutData))
		data.Copy(OutData)
		memguard.WipeBytes(OutData)
		return data, nil
	}
	cr, err = s.Rsa.MakeCrypto(s.EncrypterKeys.GetOldKey(), []byte(""))
	if err != nil {
		return nil, ErrorDecryptFile
	}
	outData2, err := cr.Decrypt(fileInfoInBytes)
	if err != nil {
		return nil, ErrorDecryptFile
	}
	data = memguard.NewBuffer(len(outData2))
	data.Copy(outData2)
	memguard.WipeBytes(outData2)
	return data, nil
}
