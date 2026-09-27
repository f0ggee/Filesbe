package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/Crypto"
	"Kaban/internal/InfrastructureLayer/FileTransferring/testFileTransfering"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"bytes"
	"crypto/rand"
	"io"
	"testing"
)

func TestNewDownloadEncrypt_decryptFile(t *testing.T) {
	key := []byte("examplekey123456") // 16 bytes key for AES-128.
	allData := make([]byte, Crypto.NewAesCtr().GetRequiredOverheadSize()+len(key))
	if _, err := io.ReadFull(rand.Reader, allData[:Crypto.NewAesCtr().GetRequiredOverheadSize()]); err != nil {
		panic(err)
	}
	copy(allData[Crypto.NewAesCtr().GetRequiredOverheadSize():], key)
	type Uploader struct {
		src  io.Reader
		dst  DomainLevel.Upload
		data []byte
	}

	type Downloader struct {
		AesKey []byte
		dst    DomainLevel.Upload
		src    DomainLevel.Download
	}
	Tests := []struct {
		name       string
		argsUpload Uploader
		argsDownl  Downloader
		wantErr    bool
	}{
		{
			name: "test1",
			argsUpload: Uploader{
				src: bytes.NewReader([]byte(rand.Text())),
				dst: func() DomainLevel.Upload {
					upload, err := testFileTransfering.NewUploader().Make(nil)
					if err != nil {
						panic(err)
					}
					return upload
				}(),
				data: func() []byte {

					return allData
				}(),
			},
			argsDownl: Downloader{
				AesKey: key,
				dst: func() DomainLevel.Upload {

					d, err := testFileTransfering.NewUpload2().Make(nil)
					if err != nil {
						panic(err)
					}
					return d
				}(),
				src: func() DomainLevel.Download {
					dow, err := testFileTransfering.NewTestDownload().Make(nil)
					if err != nil {
						panic(err)
					}
					return dow
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range Tests {
		uploaderEncrypt := &NewUploadEncrypt{
			NewUploadEncryptDataManage: NewUploadEncryptDataManage{},
			NewUploadEncryptCrypto: NewUploadEncryptCrypto{
				Generate:   nil,
				ServerKeys: DomainLevel.NewServerKeys{},
				Encrypt:    Crypto.NewAesCtr(),
			},
			NewUploadEncryptDelivery: NewUploadEncryptDelivery{},
		}
		if err := uploaderEncrypt.EncryptFile(tt.argsUpload.src, tt.argsUpload.dst, tt.argsUpload.data); err != nil && !tt.wantErr {
			t.Errorf("Encrypter: we expected = %v  but we got = %v", tt.wantErr, err)
		}

		down := &NewDownloadEncrypt{
			NewDownloadEncryptDelivery: NewDownloadEncryptDelivery{},
			NewDownloadEncryptCrypto: NewDownloadEncryptCrypto{
				Crypto:        Crypto.NewAesCtr(),
				EncrypterKeys: RepoEncrypterKeys.Keys{},
			},
			NewDownloadEncryptFileControl: NewDownloadEncryptFileControl{},
			d:                             nil,
		}

		if err := down.decryptFile(tt.argsDownl.AesKey, tt.argsDownl.src, tt.argsDownl.dst); err != nil {
			t.Errorf("Decrypter: we expected = %v  but we got = %v", tt.wantErr, err)
		}

	}

}
