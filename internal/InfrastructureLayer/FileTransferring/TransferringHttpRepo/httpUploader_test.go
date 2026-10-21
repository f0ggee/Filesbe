package TransferringHttpRepo

import (
	"bytes"
	"io"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestHttpUploader_Uploader(t *testing.T) {
	type fields struct {
		sizeFile int64
		name     string
		source   io.Writer
	}
	type args struct {
		data []byte
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
				sizeFile: 1200,
				name:     "test1",
				source:   httptest.NewRecorder(),
			},
			args: args{
				data: []byte("Hello! It's example!!!!!!!!"),
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				sizeFile: 1200,
				name:     "test2",
				source:   httptest.NewRecorder(),
			},
			args: args{
				data: []byte("Hello! IT'S EXAMPLE!!!!!!!!!!!!"),
			},
			wantErr: false,
		},
		{
			name: "test3",
			fields: fields{
				sizeFile: 1200,
				name:     "test3",
				source:   httptest.NewRecorder(),
			},
			args: args{
				data: []byte("Hello! IT'S EXAMPLE!!!!!!!!!!!!!!!!!!!"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpUploader{
				sizeFile: tt.fields.sizeFile,
				name:     tt.fields.name,
				source:   tt.fields.source,
			}
			if err := h.Upload(bytes.NewReader(tt.args.data)); (err != nil) != tt.wantErr {
				t.Errorf("Uploader() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !bytes.Equal(h.source.(*httptest.ResponseRecorder).Body.Bytes(), tt.args.data) {
				t.Errorf("Uploader() wanted = %v, got %v", h.source.(*httptest.ResponseRecorder).Body.Bytes(), tt.args.data)
			}
			if !reflect.DeepEqual(h.source.(*httptest.ResponseRecorder).Header().Get("Content-Length"), strconv.FormatInt(tt.fields.sizeFile, 10)) {
				t.Errorf("Uploader() wanted = %v, got %v", h.source.(*httptest.ResponseRecorder).Header().Get("Content-Length"), tt.fields.sizeFile)
			}
			if !strings.Contains(h.source.(*httptest.ResponseRecorder).Header().Get("Content-Disposition"), tt.fields.name) {
				t.Errorf("Uploader() wanted = %v, got %v", h.source.(*httptest.ResponseRecorder).Header().Get("Content-Disposition"), tt.fields.name)
			}

		})
	}
}
