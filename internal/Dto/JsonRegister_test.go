package Dto

import (
	"Kaban/internal/Deliver/httpController"
	"testing"
)

func TestUserDataRegister_ValidateDate(t *testing.T) {
	type fields struct {
		Name     string
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
				Name:     "Mark",
				Email:    "Mark@gmail.com",
				Password: "1234567",
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				Name:     "M",
				Email:    "Mark@gmail.com",
				Password: "1234567",
			},
			wantErr: true,
		},
		{
			name: "test3",
			fields: fields{
				Name:     "",
				Email:    "Mark@gmail.com",
				Password: "1234567",
			},
			wantErr: true,
		},
		{
			name: "test4",
			fields: fields{
				Name:     "MarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMark",
				Email:    "Mark@gmail.com",
				Password: "1234567",
			},

			wantErr: true,
		},
		{
			name: "test5",
			fields: fields{
				Name:     "Mark",
				Email:    "Markgmail.com",
				Password: "123456",
			},
			wantErr: true,
		},
		{
			name: "test6",
			fields: fields{
				Name:     "Mark",
				Email:    "MarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMarkMark@gmail.com",
				Password: "123456",
			},
			wantErr: true,
		},
		{
			name: "test7",
			fields: fields{
				Name:     "Mark",
				Email:    " ",
				Password: "123456",
			},
			wantErr: true,
		},
		{
			name: "test8",
			fields: fields{
				Name:     "Mark",
				Email:    "Mark@gmail.com",
				Password: "12345",
			},
			wantErr: true,
		},
		{
			name: "test9",
			fields: fields{
				Name:     "Mark",
				Email:    "Mark@gmail.com",
				Password: "12345123451234512345123451234512345123451234512345123451234512345123451234512345123451234512345",
			},
			wantErr: true,
		},
		{
			name: "test9",
			fields: fields{
				Name:     "Mark",
				Email:    "Mark@gmail.com",
				Password: "",
			},
			wantErr: true,
		},
		{
			name:    "test10",
			fields:  fields{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserDataRegister{
				Name:     tt.fields.Name,
				Email:    tt.fields.Email,
				Password: tt.fields.Password,
			}
			err := r.ValidateDate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr {
				t.Logf("error from the func - %v\n", err)
			}

		})
	}
}

func TestUserDataRegister_getPasswordError(t *testing.T) {
	type fields struct {
		Name     string
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
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{},
			args:    args{e: "min"},
			wantErr: true,
		},
		{
			name:    "test2",
			fields:  fields{},
			args:    args{e: "max"},
			wantErr: true,
		},
		{
			name:   "test3",
			fields: fields{},
			args: args{
				e: "required",
			},
			wantErr: true,
		},
		{
			name:   "test4",
			fields: fields{},
			args: args{
				e: "someError",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserDataRegister{
				Name:     tt.fields.Name,
				Email:    tt.fields.Email,
				Password: tt.fields.Password,
			}
			if err := r.getPasswordError(tt.args.e); (err != nil) != tt.wantErr {
				t.Errorf("getPasswordError() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr {
				t.Logf("The error from func %v\n", err)
			}
		})
	}
}

func TestUserDataRegister_getEmailError(t *testing.T) {
	type fields struct {
		Name     string
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
		wantErr bool
	}{
		{
			name:   "test1",
			fields: fields{},
			args: args{
				e: "email",
			},
			wantErr: true,
		},
		{
			name:   "test2",
			fields: fields{},
			args: args{
				e: "min",
			},
			wantErr: true,
		},
		{
			name:   "test3",
			fields: fields{},
			args: args{
				e: "max",
			},
			wantErr: true,
		},
		{
			name:   "test4",
			fields: fields{},
			args: args{
				e: "required",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserDataRegister{
				Name:     tt.fields.Name,
				Email:    tt.fields.Email,
				Password: tt.fields.Password,
			}
			if err := r.getEmailError(tt.args.e); (err != nil) != tt.wantErr {
				t.Errorf("getEmailError() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr {
				t.Logf("The errror from the func %v\n", err)
			}
		})

	}
}

func TestUserDataRegister_getNameError(t *testing.T) {
	type fields struct {
		Name     string
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
				e: "max",
			},
			wantErr: httpController.NameMaxSize,
		},
		{
			name:   "test2",
			fields: fields{},
			args: args{
				e: "min",
			},
			wantErr: httpController.NameMinSize,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserDataRegister{
				Name:     tt.fields.Name,
				Email:    tt.fields.Email,
				Password: tt.fields.Password,
			}
			if err := r.getNameError(tt.args.e); err.Error() != tt.wantErr {
				t.Errorf("the error isn't correc %v\n", err)
			}
		})
	}
}
