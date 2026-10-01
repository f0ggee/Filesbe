package Application

import (
	"Kaban/internal/DomainLevel"
	"io"
	"testing"
)

//func TestUploader(t *testing.T) {
//	key := []byte("examplekey123456")
//	// 16 bytes key for AES-128.
//	x := Crypto.NewAesCtr()
//	allData := make([]byte, x.GetRequiredOverheadSize(), x.GetRequiredOverheadSize()+len(key))
//	if _, err := io.ReadFull(rand.Reader, allData); err != nil {
//		panic(err)
//	}
//	allData = append(allData, key...)
//
//	d := NewUploadEncrypt{
//		NewUploadEncryptDataManage: NewUploadEncryptDataManage{},
//		NewUploadEncryptCrypto: NewUploadEncryptCrypto{
//			Generate:   nil,
//			ServerKeys: DomainLevel.NewServerKeys{},
//			Encrypt:    x,
//		},
//		NewUploadEncryptDelivery: NewUploadEncryptDelivery{},
//	}
//
//	testUpload1, err := testFileTransfering.NewUploader().Make(context.Background())
//	if err != nil {
//		panic(err)
//	}
//	err = EncryptFile(&d, bytes.NewReader([]byte("")), testUpload1, allData)
//	if err != nil {
//		panic(err)
//	}
//	download1, err := testFileTransfering.NewTestDownload().Make(context.Background())
//	if err != nil {
//		panic(err)
//	}
//	dow1, err := download1.Downloader()
//	if err != nil {
//		panic(err)
//	}
//	upload2, err := testFileTransfering.NewUpload2().Make(context.Background())
//	if err != nil {
//		panic(err)
//	}
//	downloadEncrypt := NewDownloadEncrypt{
//		NewDownloadEncryptDelivery: NewDownloadEncryptDelivery{},
//		NewDownloadEncryptCrypto: NewDownloadEncryptCrypto{
//			Crypto:        Crypto.NewAesCtr(),
//			EncrypterKeys: RepoEncrypterKeys.Keys{},
//		},
//		NewDownloadEncryptFileControl: NewDownloadEncryptFileControl{},
//		d:                             nil,
//	}
//
//	err = decryptFile(context.Background(), key, dow1, upload2, &downloadEncrypt)
//	if err != nil {
//		panic(err)
//	}
//}

func TestNewUploadEncrypt_EncryptFile(t *testing.T) {
	type fields struct {
		NewUploadEncryptDataManage NewUploadEncryptDataManage
		NewUploadEncryptCrypto     NewUploadEncryptCrypto
		NewUploadEncryptDelivery   NewUploadEncryptDelivery
	}
	type args struct {
		file io.Reader
		d    DomainLevel.Upload
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sa := &NewUploadEncrypt{
				NewUploadEncryptDataManage: tt.fields.NewUploadEncryptDataManage,
				NewUploadEncryptCrypto:     tt.fields.NewUploadEncryptCrypto,
				NewUploadEncryptDelivery:   tt.fields.NewUploadEncryptDelivery,
			}
			if err := sa.EncryptFile(tt.args.file, tt.args.d, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("EncryptFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
