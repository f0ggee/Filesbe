package Application

import (
	"Kaban/internal/DomainLevel"
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"log/slog"
	"time"

	"github.com/awnumar/memguard"
	"golang.org/x/sync/errgroup"
)

type NewUploadEncryptCrypto struct {
	Generate DomainLevel.CryptoGenerating
	Aes      DomainLevel.CryptoMaker
	Rsa      DomainLevel.CryptoMaker
}
type NewUploadEncryptDataManage struct {
	Encode DomainLevel.Encoder
}
type NewUploadEncryptDelivery struct {
	Uploader     DomainLevel.MakerUploader
	RedisWriter  DomainLevel.WritingRedis
	RedisChecker DomainLevel.RedisChecker
	RedisDeleter DomainLevel.DeleterRedis
	Deleter      DomainLevel.MakerDeleter
}
type UploadEncrypt struct {
	NewUploadEncryptDataManage
	NewUploadEncryptCrypto
	NewUploadEncryptDelivery
}

type UploadEncryptIncomeData struct {
	File io.Reader
	Name []byte
	Ctx  context.Context
	Size int64
}

func (u *UploadEncrypt) UploadEncryptApplication(data UploadEncryptIncomeData) (string, error) {
	defer memguard.WipeBytes(data.Name)

	if data.Size >= FileMaxSize {
		return "", ErrorFileSizeBig
	}
	g, ctx := errgroup.WithContext(data.Ctx)
	shortNameFile := u.Generate.GenerateText(4)

	x, err := u.Uploader.SetName(shortNameFile).SetSize(data.Size).Make(ctx)
	if err != nil {
		return "", err
	}
	cryptoData, err := memguard.NewBufferFromReader(rand.Reader, u.Aes.GetRequiredOverheadSize()+32)
	if err != nil {
		slog.Error("UploadEncrypt: error to make safe store", "ERROR", err)
		return "", ErrorUploadingEncrypt
	}
	defer func() {
		ctx.Done()
		_ = g.Wait()
		_ = x.CloseSource()
		cryptoData.Destroy()
	}()

	g.Go(func() error {
		err := u.encryptFile(data.File, x, cryptoData.Data(), ctx)
		if err != nil {
			return err
		}
		return nil
	})

	FileInfoInBytes, err := u.Encode.Encode(DomainLevel.FileLabelsBytes{
		FileName: data.Name,
		AesKey:   cryptoData.Data()[u.Aes.GetRequiredOverheadSize():],
	})
	if err != nil {
		return "", ErrorUploadingEncrypt
	}
	crypto, err := u.Rsa.MakeCrypto(DomainLevel.GetCryptoKeys.GetKey(), nil)
	if err != nil {
		return "", err
	}

	encryptedFileInfo, err := crypto.Encrypt(FileInfoInBytes)
	if err != nil {
		return "", err
	}

	if err = g.Wait(); err != nil {
		return "", err
	}

	err = u.RedisWriter.WriteData(DomainLevel.WriteDataIncomeData{
		FileName: shortNameFile,
		Info:     encryptedFileInfo,
		Ctx:      data.Ctx,
	})
	if err != nil {
		return "", err
	}
	u.setDeleteFile(shortNameFile)
	return shortNameFile, nil

}

func (u *UploadEncrypt) encryptFile(src io.Reader, dst DomainLevel.Upload, cryptoData []byte, ctx context.Context) error {
	if len(cryptoData) < u.Aes.GetRequiredOverheadSize() {
		return ErrorWrongCryptoData
	}
	_crypto, err := u.Aes.MakeCrypto(cryptoData[u.Aes.GetRequiredOverheadSize():], cryptoData[:u.Aes.GetRequiredOverheadSize()])
	if err != nil {
		return err
	}
	buf := memguard.NewBuffer(ChunkSize)

	defer buf.Destroy()
	for {
		select {
		case <-ctx.Done():
			return ErrorUploadingEncrypt
		default:

		}
		n, err := src.Read(buf.Bytes())
		if n > 0 {
			outData, err := _crypto.Encrypt(buf.Bytes()[:n])
			if err != nil {
				return err
			}
			err = dst.Upload(bytes.NewReader(outData))
			if err != nil {
				slog.Error("EncryptFile; error to write in the stream", "ERROR", err.Error())
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("EncryptFile; error to download a File", "ERROR", err.Error())
			return ErrorUploadingEncrypt
		}
	}
	return nil
}

func (u *UploadEncrypt) setDeleteFile(shortNameFile string) *time.Timer {
	return time.AfterFunc(5*time.Minute, func() {
		g2, Ctx := errgroup.WithContext(context.Background())
		Ctx, cancel := context.WithTimeout(Ctx, 25*time.Second)
		defer cancel()
		DownloadingHaveStarted := u.RedisChecker.ChekIsStartDownload(shortNameFile, Ctx)
		if DownloadingHaveStarted {
			return
		}
		g2.Go(func() error {

			err := u.RedisDeleter.DeleteFileInfo(shortNameFile, Ctx)
			if err != nil {
				return err
			}
			return nil
		})
		g2.Go(func() error {
			s, err := u.Deleter.SetName(shortNameFile).Make(Ctx)
			if err != nil {
				return err
			}
			return s.Deleter()
		})
		if err := g2.Wait(); err != nil {
			slog.Error("Deleter S3; error to delete a File", "ERROR", err)
			return
		}
		return
	})
}
