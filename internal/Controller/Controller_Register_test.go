package Controller

import (
	"Kaban/internal/Dto"
	"testing"
)

func TestValiDateDataForRegister(t *testing.T) {
	type args struct {
		p *Dto.UserDataRegister
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test_1",
			args: args{
				p: &Dto.UserDataRegister{
					Name:     "SELECT * FROM user",
					Email:    "SELECT * FROM USER WHERE email = user@gmail.com",
					Password: "123456",
				},
			},
			wantErr: true,
		},
		{
			name: "test_2",
			args: args{
				p: &Dto.UserDataRegister{
					Name:     "Pavel Baranov",
					Email:    "PAVELBARANOV@gmail.com",
					Password: "12345682727",
				},
			},
			wantErr: false,
		},
		{
			name: "test_3",
			args: args{
				p: &Dto.UserDataRegister{
					Name:     "SELECT * FROM user",
					Email:    "pavelBaranov@gmail.com",
					Password: "12345682727",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValiDateDataForRegister(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("ValiDateDataForRegister() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
