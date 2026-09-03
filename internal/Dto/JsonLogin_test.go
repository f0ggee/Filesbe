package Dto

import (
	"Kaban/internal/Deliver/httpController"
	"testing"
)

func TestUserLoginData_ValidateData(t *testing.T) {
	type fields struct {
		Email    string
		Password string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				Email:    "Mark@gmail.com",
				Password: "1234567",
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				Email:    "Mark@gmail.com",
				Password: "1234567123456712345671234567123456712345671234567123456712345671234567123456712345671234567123456712345671234567123456712345671234567",
			},
			wantErr: true,
		},
		{
			name: "test3",
			fields: fields{
				Email:    "Mark@gmail.com",
				Password: "1",
			},
			wantErr: true,
		},
		{
			name: "test4",
			fields: fields{
				Email:    " ",
				Password: " ",
			},

			wantErr: true,
		},
		{
			name: "test5",
			fields: fields{
				Email:    "Markgmail.com",
				Password: "123456",
			},
			wantErr: true,
		},
		{
			name: "test6",
			fields: fields{
				Email:    "MarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMark@gmail.com",
				Password: "123456",
			},
			wantErr: true,
		},
		{
			name: "test7",
			fields: fields{
				Email:    " ",
				Password: "123456",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserLoginData{
				Email:    tt.fields.Email,
				Password: tt.fields.Password,
			}
			if err := s.ValidateData(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserLoginData_getPasswordError(t *testing.T) {
	type fields struct {
		Email    string
		Password string
	}
	type args struct {
		e string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr string
	}{
		{
			name:   "test1",
			fields: fields{},
			args: args{
				e: "min",
			},
			wantErr: httpController.PasswordSizeSmall,
		},
		{
			name:   "test2",
			fields: fields{},
			args: args{
				e: "max",
			},
			wantErr: httpController.PasswordSizeBig,
		},
		{
			name:   "test3",
			fields: fields{},
			args: args{
				e: "required",
			},
			wantErr: httpController.PasswordEmpty,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserLoginData{
				Email:    tt.fields.Email,
				Password: tt.fields.Password,
			}
			if err := r.getPasswordError(tt.args.e); err.Error() != tt.wantErr {
				t.Errorf("getPasswordError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
