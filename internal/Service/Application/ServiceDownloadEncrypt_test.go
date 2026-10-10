package Application

import (
	"Kaban/internal/DomainLevel"
	"testing"
)

func TestNewDownloadEncrypt_DownloadEncrypt(t *testing.T) {
	type fields struct {
		NewDownloadEncryptDelivery    NewDownloadEncryptDelivery
		NewDownloadEncryptCrypto      NewDownloadEncryptCrypto
		NewDownloadEncryptFileControl NewDownloadEncryptFileControl
		d                             DomainLevel.Decoder
	}
	type args struct {
		data NewDownloadEncryptIncomingData
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
			s := &NewDownloadEncrypt{
				NewDownloadEncryptDelivery:    tt.fields.NewDownloadEncryptDelivery,
				NewDownloadEncryptCrypto:      tt.fields.NewDownloadEncryptCrypto,
				NewDownloadEncryptFileControl: tt.fields.NewDownloadEncryptFileControl,
				d:                             tt.fields.d,
			}
			if err := s.DownloadEncrypt(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("DownloadEncrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
