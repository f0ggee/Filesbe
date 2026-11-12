package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Drivers/Crypto"
	"Kaban/internal/Drivers/FileTransferring/s3Repo"
	"Kaban/internal/Drivers/FileTransferring/testFileTransfering"
	"Kaban/internal/Drivers/Parsers"
	"Kaban/internal/Drivers/RedisInteration"
	"Kaban/internal/Drivers/RepoEncrypterKeys"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io"
	"testing"
)

func PrepareData() ([]byte, []byte) {
	text := []byte("Hello! It's an example")
	data := make([]byte, 32)
	io.ReadFull(rand.Reader, data)
	aes, _ := Crypto.NewAesGcmEncryption().MakeCrypto(data, nil)
	cipherText, _ := aes.Encrypt(text)
	rsaK, _ := rsa.GenerateKey(rand.Reader, 2048)

	da, err := Parsers.GetNewParsing().Encode(DomainLevel.FileLabelsBytes{
		FileName: []byte("Hello!"),
		AesKey:   data,
	})
	if err != nil {
		panic(err)
	}

	DomainLevel.GetCryptoKeys = RepoEncrypterKeys.NewTestEncrypterKeys(x509.MarshalPKCS1PrivateKey(rsaK))

	rsaC, _ := Crypto.NewRsaEncryption().MakeCrypto(x509.MarshalPKCS1PrivateKey(rsaK), nil)
	labels, _ := rsaC.Encrypt(da)

	return labels, cipherText
}

func TestDownloadEncrypt_DownloadEncryptApplication(t *testing.T) {
	type fields struct {
		NewDownloadEncryptDelivery    NewDownloadEncryptDelivery
		NewDownloadEncryptCrypto      NewDownloadEncryptCrypto
		NewDownloadEncryptFileControl NewDownloadEncryptFileControl
		Decoder                       DomainLevel.Decoder
		f                             testFileTransfering.Flow
	}

	type args struct {
		data DownloadEncryptIncomeData
	}

	label, cipherText := PrepareData()
	defaultArgs := args{
		data: DownloadEncryptIncomeData{
			dst:           io.MultiWriter(),
			Ctx:           context.Background(),
			EncryptedName: "12345",
		},
	}
	defaultfields := fields{
		NewDownloadEncryptDelivery: NewDownloadEncryptDelivery{
			ReaderRedis:  RedisInteration.NewRedisTestDownloadEncrypt(label),
			DeleterRedis: RedisInteration.NewRedisDeleteTest(),
		},
		NewDownloadEncryptCrypto: NewDownloadEncryptCrypto{
			Aes: Crypto.NewAesGcmEncryption(),
			Rsa: Crypto.NewRsaEncryption(),
		},
		NewDownloadEncryptFileControl: NewDownloadEncryptFileControl{
			Uploader:   nil,
			Downloader: nil,
			Deleter:    s3Repo.NewS3DeleteTest(),
		},
		Decoder: Parsers.GetNewParsing(),
		f:       *testFileTransfering.GetTestFlow(),
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  defaultfields,
			args:    defaultArgs,
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				NewDownloadEncryptDelivery: defaultfields.NewDownloadEncryptDelivery,
				NewDownloadEncryptCrypto: NewDownloadEncryptCrypto{
					Aes: Crypto.NewAesCtrError{},
					Rsa: defaultfields.NewDownloadEncryptCrypto.Rsa,
				},
				NewDownloadEncryptFileControl: defaultfields.NewDownloadEncryptFileControl,
				Decoder:                       Parsers.GetNewParsing(),
				f:                             *testFileTransfering.GetTestFlow(),
			},
			args:    defaultArgs,
			wantErr: true,
		},
		{
			name: "test3",
			fields: fields{
				NewDownloadEncryptDelivery:    defaultfields.NewDownloadEncryptDelivery,
				NewDownloadEncryptCrypto:      defaultfields.NewDownloadEncryptCrypto,
				NewDownloadEncryptFileControl: NewDownloadEncryptFileControl{},
				Decoder:                       Parsers.GetNewParsing(),
				f:                             testFileTransfering.GetTestErrorFlow(),
			},
			args: args{
				data: DownloadEncryptIncomeData{
					dst:           io.MultiWriter(),
					Ctx:           context.Background(),
					EncryptedName: "12345",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, _ := tt.fields.f.U.Make(nil)
			x.Upload(bytes.NewReader(cipherText))
			tt.fields.NewDownloadEncryptFileControl.Uploader = tt.fields.f.U
			tt.fields.NewDownloadEncryptFileControl.Downloader = tt.fields.f.D
			s := &DownloadEncrypt{
				NewDownloadEncryptDelivery:    tt.fields.NewDownloadEncryptDelivery,
				NewDownloadEncryptCrypto:      tt.fields.NewDownloadEncryptCrypto,
				NewDownloadEncryptFileControl: tt.fields.NewDownloadEncryptFileControl,
				Decoder:                       tt.fields.Decoder,
			}
			if err := s.DownloadEncryptApplication(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("DownloadEncryptApplication() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type encryptedFileData struct {
	cipherNew       []byte
	cipherOld       []byte
	cipherUndefined []byte
}

func getFileDataPrepareData() encryptedFileData {
	d, _ := Parsers.GetNewParsing().Encode(DomainLevel.FileLabelsBytes{
		FileName: []byte("test"),
		AesKey:   []byte("test"),
	})

	rsaKeyNew, _ := rsa.GenerateKey(rand.Reader, 2048)
	rsaKeyOld, _ := rsa.GenerateKey(rand.Reader, 2048)
	rsaKeyUndefined, _ := rsa.GenerateKey(rand.Reader, 2048)

	DomainLevel.GetCryptoKeys = RepoEncrypterKeys.NewTestGetFileData()
	DomainLevel.GetCryptoKeys.UpdateOldKey(x509.MarshalPKCS1PrivateKey(rsaKeyOld))
	DomainLevel.GetCryptoKeys.UpdateNewKey(x509.MarshalPKCS1PrivateKey(rsaKeyNew))

	rsaNewKey, _ := Crypto.NewRsaEncryption().MakeCrypto(x509.MarshalPKCS1PrivateKey(rsaKeyNew), nil)
	cipherNew, _ := rsaNewKey.Encrypt(d)
	rsaOld, _ := Crypto.NewRsaEncryption().MakeCrypto(x509.MarshalPKCS1PrivateKey(rsaKeyOld), nil)
	cipherOld, _ := rsaOld.Encrypt(d)

	rsaUndefined, _ := Crypto.NewRsaEncryption().MakeCrypto(x509.MarshalPKCS1PrivateKey(rsaKeyUndefined), nil)
	cipherUndefined, _ := rsaUndefined.Encrypt(d)

	return encryptedFileData{
		cipherNew:       cipherNew,
		cipherOld:       cipherOld,
		cipherUndefined: cipherUndefined,
	}
}
func TestDownloadEncrypt_getFileData(t *testing.T) {
	data := getFileDataPrepareData()
	type fields struct {
		NewDownloadEncryptDelivery    NewDownloadEncryptDelivery
		NewDownloadEncryptCrypto      NewDownloadEncryptCrypto
		NewDownloadEncryptFileControl NewDownloadEncryptFileControl
		Decoder                       DomainLevel.Decoder
	}
	type args struct {
		fileInfoInBytes []byte
	}
	defaultFields := fields{
		NewDownloadEncryptCrypto: NewDownloadEncryptCrypto{
			Rsa: Crypto.NewRsaEncryption(),
		},
		Decoder: Parsers.GetNewParsing(),
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test1",
			fields: defaultFields,
			args: args{
				fileInfoInBytes: data.cipherNew,
			},
			wantErr: false,
		},
		{
			name:   "test2",
			fields: defaultFields,
			args: args{
				fileInfoInBytes: data.cipherOld,
			},
			wantErr: false,
		},
		{
			name:   "test3",
			fields: defaultFields,
			args: args{
				fileInfoInBytes: data.cipherUndefined,
			},
			wantErr: true,
		},
		{
			name: "test4",
			fields: fields{
				NewDownloadEncryptCrypto: NewDownloadEncryptCrypto{
					Rsa: Crypto.NewRsaError(),
				},
			},
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DownloadEncrypt{
				NewDownloadEncryptDelivery:    tt.fields.NewDownloadEncryptDelivery,
				NewDownloadEncryptCrypto:      tt.fields.NewDownloadEncryptCrypto,
				NewDownloadEncryptFileControl: tt.fields.NewDownloadEncryptFileControl,
				Decoder:                       tt.fields.Decoder,
			}
			d1, d2, err := s.getFileData(tt.args.fileInfoInBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if d1 == nil != tt.wantErr {
				t.Errorf("getFileData() a first key is nil %v", d1)
				return
			}
			if d2 == nil != tt.wantErr {
				t.Errorf("getFileData() a first key is nil %v", d2)
			}
		})
	}
}

func TestDownloadEncrypt_decryptFile(t *testing.T) {

	prepareDataNormal := func(inputData []byte) ([]byte, []byte) {

		key := make([]byte, 32)
		io.ReadFull(rand.Reader, key)
		_crypto, _ := Crypto.NewAesGcmEncryption().MakeCrypto(key, nil)
		cipherText, _ := _crypto.Encrypt(key)

		return cipherText, key
	}
	prepareDataError := func(inputData []byte) ([]byte, []byte) {
		key := make([]byte, 32)
		io.ReadFull(rand.Reader, key)
		_crypto, _ := Crypto.NewAesGcmEncryption().MakeCrypto(key, nil)
		cipherText, _ := _crypto.Encrypt(key)
		return cipherText, make([]byte, 32)
	}
	type fields struct {
		NewDownloadEncryptCrypto NewDownloadEncryptCrypto
	}
	type args struct {
		AesKey []byte
		src    DomainLevel.Download
		dst    DomainLevel.Upload
	}
	defaultFields := fields{
		NewDownloadEncryptCrypto: NewDownloadEncryptCrypto{
			Aes: Crypto.NewAesGcmEncryption(),
		},
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		inputData []byte
		flow      testFileTransfering.TestDecryptFile
		data      func(inputData []byte) ([]byte, []byte)
		wantErr   bool
	}{
		{
			name:      "test1",
			fields:    defaultFields,
			args:      args{},
			inputData: []byte("Hello!Hello!Hello!Hello!Hello!Hello!Hello!Hello!Hello!Hello!Hello!Hello!Hello!Hello!Hello!Hello!"),
			flow:      testFileTransfering.GetNormalFlow(),
			data:      prepareDataNormal,
			wantErr:   false,
		},
		{
			name:      "test2",
			fields:    defaultFields,
			args:      args{},
			inputData: []byte("Hello!"),
			flow:      testFileTransfering.GetNormalFlow(),
			data:      prepareDataError,
			wantErr:   true,
		},
		{
			name:      "test3",
			fields:    defaultFields,
			args:      args{},
			inputData: []byte("TestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTest"),
			flow:      testFileTransfering.GetTestErrorFlowErrorDecryptFile(),
			data:      prepareDataNormal,
			wantErr:   true,
		},
		{
			name: "test4",
			fields: fields{
				NewDownloadEncryptCrypto: NewDownloadEncryptCrypto{
					Aes: Crypto.NewAesCtrError{},
				},
			},
			args:      args{},
			inputData: []byte("TestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTestTest"),
			flow:      testFileTransfering.GetNormalFlow(),
			data:      prepareDataNormal,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DownloadEncrypt{
				NewDownloadEncryptCrypto: tt.fields.NewDownloadEncryptCrypto,
			}
			cipherText, key := tt.data(tt.inputData)

			tt.flow.SetCipherText(cipherText)
			tt.args.AesKey = key
			_u, _ := tt.flow.U.Make(nil)
			tt.args.dst = _u
			_d, _ := tt.flow.D.Make(nil)
			tt.args.src = _d
			if err := s.decryptFile(tt.args.AesKey, tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("decryptFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.flow.Reset()

		})
	}
}
