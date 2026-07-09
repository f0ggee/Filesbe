package RepoEncrypterKeys

import (
	"crypto/rand"
	"testing"

	"github.com/awnumar/memguard"
)

func TestKeys_UpdateNewKey(t *testing.T) {
	type fields struct {
		OldKey         *memguard.LockedBuffer
		NewKey         *memguard.LockedBuffer
		isOldKeyCreate bool
	}

	newKey, _ := memguard.NewBufferFromReader(rand.Reader, 32)
	type args struct {
		key *memguard.LockedBuffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				OldKey:         &memguard.LockedBuffer{},
				NewKey:         &memguard.LockedBuffer{},
				isOldKeyCreate: false,
			},
			wantErr: true,
		},
		{
			name: "test2",
			fields: fields{
				OldKey:         &memguard.LockedBuffer{},
				NewKey:         newKey,
				isOldKeyCreate: true,
			},
			args:    args{key: newKey},
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
			if err := s.UpdateNewKey(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("UpdateNewKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
