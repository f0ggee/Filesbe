package RepoEncrypterKeys

import (
	"crypto/rand"
	"testing"

	"github.com/awnumar/memguard"
)

func GetMem() *memguard.LockedBuffer {

	d, err := memguard.NewBufferFromReader(rand.Reader, 32)
	if err != nil {
		panic(err)
	}
	return d
}

func TestKeys_UpdateNewKey(t *testing.T) {
	type fields struct {
		OldKey         *memguard.LockedBuffer
		NewKey         *memguard.LockedBuffer
		isOldKeyCreate bool
	}
	type args struct {
		key []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				OldKey:         nil,
				NewKey:         nil,
				isOldKeyCreate: false,
			},
			args: args{
				key: []byte(rand.Text()),
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				OldKey:         nil,
				NewKey:         GetMem(),
				isOldKeyCreate: false,
			},
			args: args{
				key: []byte(rand.Text()),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Keys{
				OldKey:         tt.fields.OldKey,
				NewKey:         tt.fields.NewKey,
				isOldKeyCreate: tt.fields.isOldKeyCreate,
			}
			s.UpdateNewKey(tt.args.key)
		})
	}
}

func TestKeys_UpdateOldKey(t *testing.T) {
	type fields struct {
		OldKey         *memguard.LockedBuffer
		NewKey         *memguard.LockedBuffer
		isOldKeyCreate bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "test1",
			fields: fields{
				OldKey:         nil,
				NewKey:         nil,
				isOldKeyCreate: false,
			},
		},
		{
			name: "test2",
			fields: fields{
				OldKey:         nil,
				NewKey:         GetMem(),
				isOldKeyCreate: false,
			},
		},
		{
			name: "test3",
			fields: fields{
				OldKey:         GetMem(),
				NewKey:         GetMem(),
				isOldKeyCreate: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Keys{
				OldKey:         tt.fields.OldKey,
				NewKey:         tt.fields.NewKey,
				isOldKeyCreate: tt.fields.isOldKeyCreate,
			}
			s.UpdateOldKey()
		})
	}
}
