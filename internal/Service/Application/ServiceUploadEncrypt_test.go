package Application

import "testing"

func TestUploadEncrypt_UploadEncryptApplication(t *testing.T) {
	type fields struct {
		NewUploadEncryptDataManage NewUploadEncryptDataManage
		NewUploadEncryptCrypto     NewUploadEncryptCrypto
		NewUploadEncryptDelivery   NewUploadEncryptDelivery
	}
	type args struct {
		data UploadEncryptIncomeData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sa := &UploadEncrypt{
				NewUploadEncryptDataManage: tt.fields.NewUploadEncryptDataManage,
				NewUploadEncryptCrypto:     tt.fields.NewUploadEncryptCrypto,
				NewUploadEncryptDelivery:   tt.fields.NewUploadEncryptDelivery,
			}
			got, err := sa.UploadEncryptApplication(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadEncryptApplication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UploadEncryptApplication() got = %v, want %v", got, tt.want)
			}
		})
	}
}
