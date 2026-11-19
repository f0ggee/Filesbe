package httpController

import (
	"Kaban/internal/Drivers/Tokens"
	"Kaban/internal/Service/Application"
	"Kaban/internal/Service/Application/testFunctions"
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadController_FileUploader(t *testing.T) {
	type fields struct {
		FileUploaderSessions FileUploaderSessions
		UploadApplication    Application.UploadApplication
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantError bool
	}{
		{
			name: "test1",
			fields: fields{
				FileUploaderSessions: FileUploaderSessions{
					Rf:       Tokens.NewRfToken(),
					Jwt:      Tokens.NewJwtToken(),
					Sessions: newTestSession(),
				},
				UploadApplication: testFunctions.TestUpload(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: newFileUploadRequest("/data", "File", "example.txt", []byte("Hello!")),
			},
			wantError: false,
		},
		{
			name: "test2",
			fields: fields{
				FileUploaderSessions: FileUploaderSessions{
					Rf:       Tokens.NewRfToken(),
					Jwt:      Tokens.NewJwtToken(),
					Sessions: newTestBadSession(),
				},
				UploadApplication: testFunctions.TestUpload(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: newFileUploadRequest("/data", "File", "example.txt", []byte("Hello!")),
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &UploadController{
				FileUploaderSessions: tt.fields.FileUploaderSessions,
				UploadApplication:    tt.fields.UploadApplication,
			}
			d.FileUploader(tt.args.w, tt.args.r)

			var answer AnswerUploaderFileNoEncrypt
			if err := json.NewDecoder(tt.args.w.Body).Decode(&answer); err != nil {
				t.Fatalf("FileUploader() can't parse the answer: %v", err)
			}
			if tt.wantError && answer.Error == "" {
				t.Errorf("FileUploader() expected an error, got none, answer = %+v", answer)
			}
			if !tt.wantError && answer.Error != "" {
				t.Errorf("FileUploader() unexpected error = %s", answer.Error)
			}
		})
	}
}
func newFileUploadRequest(url, fieldName, fileName string, content []byte) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile(fieldName, fileName)
	part.Write(content)
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}
