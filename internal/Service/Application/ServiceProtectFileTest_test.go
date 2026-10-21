package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/Crypto"
	"Kaban/internal/InfrastructureLayer/FileTransferring/testFileTransfering"
	"bytes"
	"crypto/rand"
	"io"
	"testing"
)

func TestProtectFile_protectFile(t *testing.T) {

	crypt := Crypto.NewAesGcmEncryption()
	cryptoData := make([]byte, 32+crypt.GetRequiredOverheadSize())
	if _, err := io.ReadFull(rand.Reader, cryptoData); err != nil {
		panic(err)
	}
	firstUpload := testFileTransfering.GetFirstUploading()
	second := testFileTransfering.GetUploaderSecond()
	secondUpload, err := second.Make(nil)
	if err != nil {
		panic(err)
	}
	down := testFileTransfering.GetTestDownloader()
	type fields struct {
		Encrypt DomainLevel.CryptoMaker
	}
	type args struct {
		src  io.Reader
		dst  DomainLevel.Upload
		data []byte
	}

	type fieldsDownloadCrypto struct {
		Aes DomainLevel.CryptoMaker
	}
	type argsDownloadCrypto struct {
		AesKey      []byte
		o           DomainLevel.Download
		uploadedObj DomainLevel.Upload
	}
	tests := []struct {
		name           string
		fieldsUpload   fields
		argsUpload     args
		fieldsDownload fieldsDownloadCrypto
		argsDownload   argsDownloadCrypto
		wantErr        bool
	}{
		{
			name: "test1",
			fieldsUpload: fields{
				Encrypt: crypt,
			},
			argsUpload: args{
				src:  bytes.NewReader([]byte("Hello! It's example")),
				dst:  &firstUpload,
				data: cryptoData,
			},
			fieldsDownload: fieldsDownloadCrypto{
				Aes: crypt,
			},
			argsDownload: argsDownloadCrypto{
				AesKey:      cryptoData[crypt.GetRequiredOverheadSize():],
				o:           down,
				uploadedObj: secondUpload,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sa := &EncryptFile{
				Encrypt: tt.fieldsUpload.Encrypt,
			}
			if err := sa.encryptFile(tt.argsUpload.src, tt.argsUpload.dst, tt.argsUpload.data); (err != nil) != tt.wantErr {
				t.Errorf("encryptFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			da := &DecryptFile{
				Aes: tt.fieldsDownload.Aes,
			}

			if err := da.decryptFile(tt.argsDownload.AesKey[crypt.GetRequiredOverheadSize():], tt.argsDownload.o, tt.argsDownload.uploadedObj); err != nil {
				t.Errorf("decryptFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
